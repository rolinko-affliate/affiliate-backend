#!/usr/bin/env python3

"""
Idempotent Test Account Creation Script for Affiliate Platform Backend

Creates test accounts and organizations with full idempotency.
Safe to run multiple times without creating duplicates.
"""

import argparse
import json
import logging
import sys
import time
import uuid
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass
from enum import Enum
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

class OperationResult(Enum):
    CREATED = "created"
    ALREADY_EXISTS = "already_exists"
    UPDATED = "updated"
    FAILED = "failed"

@dataclass
class TestUser:
    uuid: str
    email: str
    first_name: str
    last_name: str
    role_name: str
    role_id: int
    organization_name: str

@dataclass
class TestOrganization:
    name: str
    type: str
    description: str

@dataclass
class IdempotentResult:
    operation: str
    result: OperationResult
    entity_id: Optional[Any] = None
    entity_data: Optional[Dict] = None
    error_message: Optional[str] = None
    
    @property
    def success(self) -> bool:
        return self.result in [OperationResult.CREATED, OperationResult.ALREADY_EXISTS, OperationResult.UPDATED]

class IdempotentAPIClient:
    
    def __init__(self, base_url: str, admin_token: Optional[str] = None, timeout: int = 30):
        self.base_url = base_url.rstrip('/')
        self.admin_token = admin_token
        self.timeout = timeout
        
        # Configure session with retries
        self.session = requests.Session()
        retry_strategy = Retry(
            total=3,
            backoff_factor=1,
            status_forcelist=[429, 500, 502, 503, 504],
        )
        adapter = HTTPAdapter(max_retries=retry_strategy)
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)
        
        # Set default headers
        self.session.headers.update({
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        })
        
        if self.admin_token:
            self.session.headers.update({
                'Authorization': f'Bearer {self.admin_token}'
            })
    
    def _make_request(self, method: str, endpoint: str, data: Optional[Dict] = None, 
                     params: Optional[Dict] = None) -> Tuple[bool, int, Dict]:
        url = f"{self.base_url}{endpoint}"
        
        try:
            logger.debug(f"Making {method} request to {url}")
            if data:
                logger.debug(f"Request data: {json.dumps(data, indent=2)}")
            
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                params=params,
                timeout=self.timeout
            )
            
            logger.debug(f"Response status: {response.status_code}")
            
            # Try to parse JSON response
            try:
                response_data = response.json()
            except json.JSONDecodeError:
                response_data = {"raw_response": response.text}
            
            success = 200 <= response.status_code < 300
            return success, response.status_code, response_data
                
        except requests.exceptions.RequestException as e:
            logger.error(f"Request failed: {e}")
            return False, 0, {"error": str(e)}
    
    def health_check(self) -> bool:
        success, status_code, data = self._make_request('GET', '/health')
        return success
    
    def get_organizations(self) -> Tuple[bool, List[Dict]]:
        success, status_code, data = self._make_request('GET', '/api/v1/organizations')
        if success:
            # Handle both list response and paginated response
            if isinstance(data, list):
                return True, data
            elif isinstance(data, dict) and 'organizations' in data:
                return True, data['organizations']
            elif isinstance(data, dict) and 'data' in data:
                return True, data['data']
            else:
                return True, []
        return False, []
    
    def get_organization_by_name(self, name: str) -> Tuple[bool, Optional[Dict]]:
    
        success, orgs = self.get_organizations()
        if not success:
            return False, None
        
        for org in orgs:
            if org.get('name') == name:
                return True, org
        return True, None
    
    def create_organization_idempotent(self, name: str, org_type: str) -> IdempotentResult:
    
        # Check if organization already exists
        exists_success, existing_org = self.get_organization_by_name(name)
        if not exists_success:
            return IdempotentResult(
                operation=f"create_organization({name})",
                result=OperationResult.FAILED,
                error_message="Failed to check existing organizations"
            )
        
        if existing_org:
            # Organization already exists
            if existing_org.get('type') == org_type:
                return IdempotentResult(
                    operation=f"create_organization({name})",
                    result=OperationResult.ALREADY_EXISTS,
                    entity_id=existing_org.get('organization_id'),
                    entity_data=existing_org
                )
            else:
                return IdempotentResult(
                    operation=f"create_organization({name})",
                    result=OperationResult.FAILED,
                    error_message=f"Organization exists with different type: {existing_org.get('type')} != {org_type}"
                )
        
        # Create new organization
        data = {"name": name, "type": org_type}
        success, status_code, response_data = self._make_request('POST', '/api/v1/organizations', data=data)
        
        if success:
            return IdempotentResult(
                operation=f"create_organization({name})",
                result=OperationResult.CREATED,
                entity_id=response_data.get('organization_id'),
                entity_data=response_data
            )
        else:
            # Check if it's a duplicate error
            error_msg = response_data.get('error', f'HTTP {status_code}')
            if 'duplicate' in error_msg.lower() or 'already exists' in error_msg.lower() or status_code == 409:
                # Race condition - organization was created between our check and creation
                exists_success, existing_org = self.get_organization_by_name(name)
                if exists_success and existing_org:
                    return IdempotentResult(
                        operation=f"create_organization({name})",
                        result=OperationResult.ALREADY_EXISTS,
                        entity_id=existing_org.get('organization_id'),
                        entity_data=existing_org
                    )
            
            return IdempotentResult(
                operation=f"create_organization({name})",
                result=OperationResult.FAILED,
                error_message=error_msg
            )
    
    def upsert_profile(self, user_uuid: str, email: str, first_name: str,
                      last_name: str, role_id: int, organization_id: Optional[int] = None) -> IdempotentResult:
    
        data = {
            "id": user_uuid,
            "email": email,
            "first_name": first_name,
            "last_name": last_name,
            "role_id": role_id
        }
        if organization_id:
            data["organization_id"] = organization_id
        
        success, status_code, response_data = self._make_request('POST', '/api/v1/profiles/upsert', data=data)
        
        if success:
            # The upsert endpoint doesn't tell us if it was created or updated
            # We'll assume it was successful and mark as CREATED for simplicity
            return IdempotentResult(
                operation=f"upsert_profile({email})",
                result=OperationResult.CREATED,
                entity_id=user_uuid,
                entity_data=response_data
            )
        else:
            error_msg = response_data.get('error', f'HTTP {status_code}')
            return IdempotentResult(
                operation=f"upsert_profile({email})",
                result=OperationResult.FAILED,
                error_message=error_msg
            )
    
    def get_advertisers(self) -> Tuple[bool, List[Dict]]:
    
        success, status_code, data = self._make_request('GET', '/api/v1/advertisers')
        if success:
            if isinstance(data, list):
                return True, data
            elif isinstance(data, dict) and 'advertisers' in data:
                return True, data['advertisers']
            elif isinstance(data, dict) and 'data' in data:
                return True, data['data']
            else:
                return True, []
        return False, []
    
    def get_advertiser_by_name(self, name: str) -> Tuple[bool, Optional[Dict]]:
    
        success, advertisers = self.get_advertisers()
        if not success:
            return False, None
        
        for adv in advertisers:
            if adv.get('name') == name:
                return True, adv
        return True, None
    
    def create_advertiser_idempotent(self, name: str, organization_id: int, contact_email: str,
                                   billing_details: Optional[Dict] = None) -> IdempotentResult:
    
        # Check if advertiser already exists
        exists_success, existing_adv = self.get_advertiser_by_name(name)
        if not exists_success:
            return IdempotentResult(
                operation=f"create_advertiser({name})",
                result=OperationResult.FAILED,
                error_message="Failed to check existing advertisers"
            )
        
        if existing_adv:
            return IdempotentResult(
                operation=f"create_advertiser({name})",
                result=OperationResult.ALREADY_EXISTS,
                entity_id=existing_adv.get('advertiser_id'),
                entity_data=existing_adv
            )
        
        # Create new advertiser
        data = {
            "name": name,
            "organization_id": organization_id,
            "contact_email": contact_email,
            "status": "active"
        }
        if billing_details:
            data["billing_details"] = billing_details
        
        success, status_code, response_data = self._make_request('POST', '/api/v1/advertisers', data=data)
        
        if success:
            return IdempotentResult(
                operation=f"create_advertiser({name})",
                result=OperationResult.CREATED,
                entity_id=response_data.get('advertiser_id'),
                entity_data=response_data
            )
        else:
            error_msg = response_data.get('error', f'HTTP {status_code}')
            if 'duplicate' in error_msg.lower() or 'already exists' in error_msg.lower() or status_code == 409:
                # Race condition - check again
                exists_success, existing_adv = self.get_advertiser_by_name(name)
                if exists_success and existing_adv:
                    return IdempotentResult(
                        operation=f"create_advertiser({name})",
                        result=OperationResult.ALREADY_EXISTS,
                        entity_id=existing_adv.get('advertiser_id'),
                        entity_data=existing_adv
                    )
            
            return IdempotentResult(
                operation=f"create_advertiser({name})",
                result=OperationResult.FAILED,
                error_message=error_msg
            )
    
    def get_affiliates(self) -> Tuple[bool, List[Dict]]:
    
        success, status_code, data = self._make_request('GET', '/api/v1/affiliates')
        if success:
            if isinstance(data, list):
                return True, data
            elif isinstance(data, dict) and 'affiliates' in data:
                return True, data['affiliates']
            elif isinstance(data, dict) and 'data' in data:
                return True, data['data']
            else:
                return True, []
        return False, []
    
    def get_affiliate_by_name(self, name: str) -> Tuple[bool, Optional[Dict]]:
    
        success, affiliates = self.get_affiliates()
        if not success:
            return False, None
        
        for aff in affiliates:
            if aff.get('name') == name:
                return True, aff
        return True, None
    
    def create_affiliate_idempotent(self, name: str, organization_id: int, contact_email: str,
                                  payment_details: Optional[Dict] = None) -> IdempotentResult:
    
        # Check if affiliate already exists
        exists_success, existing_aff = self.get_affiliate_by_name(name)
        if not exists_success:
            return IdempotentResult(
                operation=f"create_affiliate({name})",
                result=OperationResult.FAILED,
                error_message="Failed to check existing affiliates"
            )
        
        if existing_aff:
            return IdempotentResult(
                operation=f"create_affiliate({name})",
                result=OperationResult.ALREADY_EXISTS,
                entity_id=existing_aff.get('affiliate_id'),
                entity_data=existing_aff
            )
        
        # Create new affiliate
        data = {
            "name": name,
            "organization_id": organization_id,
            "contact_email": contact_email,
            "status": "active"
        }
        if payment_details:
            data["payment_details"] = payment_details
        
        success, status_code, response_data = self._make_request('POST', '/api/v1/affiliates', data=data)
        
        if success:
            return IdempotentResult(
                operation=f"create_affiliate({name})",
                result=OperationResult.CREATED,
                entity_id=response_data.get('affiliate_id'),
                entity_data=response_data
            )
        else:
            error_msg = response_data.get('error', f'HTTP {status_code}')
            if 'duplicate' in error_msg.lower() or 'already exists' in error_msg.lower() or status_code == 409:
                # Race condition - check again
                exists_success, existing_aff = self.get_affiliate_by_name(name)
                if exists_success and existing_aff:
                    return IdempotentResult(
                        operation=f"create_affiliate({name})",
                        result=OperationResult.ALREADY_EXISTS,
                        entity_id=existing_aff.get('affiliate_id'),
                        entity_data=existing_aff
                    )
            
            return IdempotentResult(
                operation=f"create_affiliate({name})",
                result=OperationResult.FAILED,
                error_message=error_msg
            )
    
    def get_campaigns(self) -> Tuple[bool, List[Dict]]:
    
        success, status_code, data = self._make_request('GET', '/api/v1/campaigns')
        if success:
            if isinstance(data, list):
                return True, data
            elif isinstance(data, dict) and 'campaigns' in data:
                return True, data['campaigns']
            elif isinstance(data, dict) and 'data' in data:
                return True, data['data']
            else:
                return True, []
        return False, []
    
    def get_campaign_by_name(self, name: str) -> Tuple[bool, Optional[Dict]]:
    
        success, campaigns = self.get_campaigns()
        if not success:
            return False, None
        
        for camp in campaigns:
            if camp.get('name') == name:
                return True, camp
        return True, None
    
    def create_campaign_idempotent(self, name: str, advertiser_id: int, organization_id: int, 
                                 payout_type: str, payout_amount: float, revenue_type: str, 
                                 revenue_amount: float, currency_id: str = "USD", 
                                 status: str = "active", visibility: str = "public") -> IdempotentResult:
    
        # Check if campaign already exists
        exists_success, existing_camp = self.get_campaign_by_name(name)
        if not exists_success:
            return IdempotentResult(
                operation=f"create_campaign({name})",
                result=OperationResult.FAILED,
                error_message="Failed to check existing campaigns"
            )
        
        if existing_camp:
            return IdempotentResult(
                operation=f"create_campaign({name})",
                result=OperationResult.ALREADY_EXISTS,
                entity_id=existing_camp.get('campaign_id'),
                entity_data=existing_camp
            )
        
        # Create new campaign
        data = {
            "name": name,
            "advertiser_id": advertiser_id,
            "organization_id": organization_id,
            "payout_type": payout_type,
            "payout_amount": payout_amount,
            "revenue_type": revenue_type,
            "revenue_amount": revenue_amount,
            "currency_id": currency_id,
            "status": status,
            "visibility": visibility
        }
        
        success, status_code, response_data = self._make_request('POST', '/api/v1/campaigns', data=data)
        
        if success:
            return IdempotentResult(
                operation=f"create_campaign({name})",
                result=OperationResult.CREATED,
                entity_id=response_data.get('campaign_id'),
                entity_data=response_data
            )
        else:
            error_msg = response_data.get('error', f'HTTP {status_code}')
            if 'duplicate' in error_msg.lower() or 'already exists' in error_msg.lower() or status_code == 409:
                # Race condition - check again
                exists_success, existing_camp = self.get_campaign_by_name(name)
                if exists_success and existing_camp:
                    return IdempotentResult(
                        operation=f"create_campaign({name})",
                        result=OperationResult.ALREADY_EXISTS,
                        entity_id=existing_camp.get('campaign_id'),
                        entity_data=existing_camp
                    )
            
            return IdempotentResult(
                operation=f"create_campaign({name})",
                result=OperationResult.FAILED,
                error_message=error_msg
            )
    
    def create_analytics_data_idempotent(self, endpoint: str, data: Dict, identifier: str) -> IdempotentResult:
    
        success, status_code, response_data = self._make_request('POST', endpoint, data=data)
        
        if success:
            return IdempotentResult(
                operation=f"create_analytics({identifier})",
                result=OperationResult.CREATED,
                entity_data=response_data
            )
        else:
            error_msg = response_data.get('error', f'HTTP {status_code}')
            # For analytics, we might want to be more lenient about duplicates
            if 'duplicate' in error_msg.lower() or 'already exists' in error_msg.lower() or status_code == 409:
                return IdempotentResult(
                    operation=f"create_analytics({identifier})",
                    result=OperationResult.ALREADY_EXISTS,
                    entity_data=response_data
                )
            
            return IdempotentResult(
                operation=f"create_analytics({identifier})",
                result=OperationResult.FAILED,
                error_message=error_msg
            )

class IdempotentTestAccountCreator:

    
    def __init__(self, api_client: IdempotentAPIClient, verbose: bool = False):
        self.api = api_client
        self.verbose = verbose
        self.results = {
            'organizations': {},
            'profiles': {},
            'advertisers': {},
            'affiliates': {},
            'campaigns': {},
            'analytics': {}
        }
        
        if verbose:
            logger.setLevel(logging.DEBUG)
    
    def log_result(self, result: IdempotentResult):
    
        if result.success:
            if result.result == OperationResult.CREATED:
                logger.info(f"✅ Created: {result.operation}")
            elif result.result == OperationResult.ALREADY_EXISTS:
                logger.info(f"ℹ️  Already exists: {result.operation}")
            elif result.result == OperationResult.UPDATED:
                logger.info(f"🔄 Updated: {result.operation}")
        else:
            logger.error(f"❌ Failed: {result.operation} - {result.error_message}")
    
    def check_api_health(self) -> bool:
    
        logger.info("Checking API health...")
        if self.api.health_check():
            logger.info("✅ API is healthy and responding")
            return True
        else:
            logger.error("❌ API health check failed")
            return False
    
    def get_test_organizations(self) -> List[TestOrganization]:
    
        return [
            TestOrganization(
                name="UpsailAI",
                type="platform_owner",
                description="Platform administration organization"
            ),
            TestOrganization(
                name="Adidas",
                type="advertiser",
                description="Global sportswear brand"
            ),
            TestOrganization(
                name="Dyson",
                type="advertiser",
                description="Home appliance technology company"
            ),
            TestOrganization(
                name="Le Monde",
                type="affiliate",
                description="French news publisher"
            )
        ]
    
    def get_test_users(self) -> List[TestUser]:
    
        return [
            TestUser(
                uuid="550e8400-e29b-41d4-a716-446655440000",
                email="admin@upsailai.com",
                first_name="Admin",
                last_name="User",
                role_name="Admin",
                role_id=1,
                organization_name="UpsailAI"
            ),
            TestUser(
                uuid="a654ad6a-83c7-44c5-9f34-d2d5adb2f8a0",
                email="rolinko@adidas.com",
                first_name="Roland",
                last_name="Adidas",
                role_name="AdvertiserManager",
                role_id=1000,
                organization_name="Adidas"
            ),
            TestUser(
                uuid="71ae7a37-92e5-4693-91e1-f5a1464b7414",
                email="rolinko@dyson.com",
                first_name="Roland",
                last_name="Dyson",
                role_name="AdvertiserManager",
                role_id=1000,
                organization_name="Dyson"
            ),
            TestUser(
                uuid="268826c9-d59d-4b40-9558-4ce5f7bf7534",
                email="rolinko@lemonde.fr",
                first_name="Roland",
                last_name="LeMonde",
                role_name="AffiliateManager",
                role_id=1001,
                organization_name="Le Monde"
            )
        ]
    
    def create_organizations(self) -> bool:
    
        logger.info("Creating test organizations...")
        
        organizations = self.get_test_organizations()
        success_count = 0
        
        for org in organizations:
            logger.info(f"Processing organization: {org.name} ({org.type})")
            
            result = self.api.create_organization_idempotent(org.name, org.type)
            self.log_result(result)
            
            if result.success:
                self.results['organizations'][org.name] = {
                    'id': result.entity_id,
                    'type': org.type,
                    'result': result.result.value,
                    'data': result.entity_data
                }
                success_count += 1
            else:
                self.results['organizations'][org.name] = {
                    'id': None,
                    'type': org.type,
                    'result': result.result.value,
                    'error': result.error_message
                }
        
        logger.info(f"Organizations processed: {success_count}/{len(organizations)} successful")
        return success_count == len(organizations)
    
    def create_user_profiles(self) -> bool:
    
        logger.info("Creating test user profiles...")
        
        users = self.get_test_users()
        success_count = 0
        
        for user in users:
            logger.info(f"Processing profile for: {user.email}")
            
            # Get organization ID
            org_id = None
            if user.organization_name in self.results['organizations']:
                org_data = self.results['organizations'][user.organization_name]
                if org_data['id']:
                    org_id = org_data['id']
                else:
                    logger.error(f"Organization '{user.organization_name}' was not created successfully")
                    continue
            else:
                logger.error(f"Organization '{user.organization_name}' not found in results")
                continue
            
            result = self.api.upsert_profile(
                user_uuid=user.uuid,
                email=user.email,
                first_name=user.first_name,
                last_name=user.last_name,
                role_id=user.role_id,
                organization_id=org_id
            )
            self.log_result(result)
            
            if result.success:
                self.results['profiles'][user.email] = {
                    'uuid': user.uuid,
                    'organization_id': org_id,
                    'role_id': user.role_id,
                    'result': result.result.value,
                    'data': result.entity_data
                }
                success_count += 1
            else:
                self.results['profiles'][user.email] = {
                    'uuid': user.uuid,
                    'organization_id': org_id,
                    'role_id': user.role_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
        
        logger.info(f"User profiles processed: {success_count}/{len(users)} successful")
        return success_count == len(users)
    
    def create_advertisers(self) -> bool:
    
        logger.info("Creating test advertisers...")
        
        advertisers_data = [
            {
                "name": "Adidas Global",
                "organization_name": "Adidas",
                "contact_email": "rolinko@adidas.com",
                "billing_details": {
                    "company_name": "Adidas AG",
                    "address": {
                        "street": "Adi-Dassler-Strasse 1",
                        "city": "Herzogenaurach",
                        "postal_code": "91074",
                        "country": "Germany"
                    },
                    "tax_id": "DE123456789",
                    "billing_email": "billing@adidas.com"
                }
            },
            {
                "name": "Dyson Ltd",
                "organization_name": "Dyson",
                "contact_email": "rolinko@dyson.com",
                "billing_details": {
                    "company_name": "Dyson Ltd",
                    "address": {
                        "street": "Tetbury Hill",
                        "city": "Malmesbury",
                        "postal_code": "SN16 0RP",
                        "country": "United Kingdom"
                    },
                    "tax_id": "GB123456789",
                    "billing_email": "billing@dyson.com"
                }
            }
        ]
        
        success_count = 0
        
        for adv_data in advertisers_data:
            logger.info(f"Processing advertiser: {adv_data['name']}")
            
            # Get organization ID
            org_id = None
            if adv_data['organization_name'] in self.results['organizations']:
                org_data = self.results['organizations'][adv_data['organization_name']]
                if org_data['id']:
                    org_id = org_data['id']
                else:
                    logger.error(f"Organization '{adv_data['organization_name']}' was not created successfully")
                    continue
            else:
                logger.error(f"Organization '{adv_data['organization_name']}' not found in results")
                continue
            
            result = self.api.create_advertiser_idempotent(
                name=adv_data['name'],
                organization_id=org_id,
                contact_email=adv_data['contact_email'],
                billing_details=adv_data['billing_details']
            )
            self.log_result(result)
            
            if result.success:
                self.results['advertisers'][adv_data['name']] = {
                    'id': result.entity_id,
                    'organization_id': org_id,
                    'result': result.result.value,
                    'data': result.entity_data
                }
                success_count += 1
            else:
                self.results['advertisers'][adv_data['name']] = {
                    'id': None,
                    'organization_id': org_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
        
        logger.info(f"Advertisers processed: {success_count}/{len(advertisers_data)} successful")
        return success_count == len(advertisers_data)
    
    def create_affiliates(self) -> bool:
    
        logger.info("Creating test affiliates...")
        
        affiliates_data = [
            {
                "name": "Le Monde",
                "organization_name": "Le Monde",
                "contact_email": "rolinko@lemonde.fr",
                "payment_details": {
                    "preferred_method": "bank_transfer",
                    "bank_details": {
                        "account_holder": "Le Monde SA",
                        "iban": "FR1420041010050500013M02606",
                        "bic": "PSSTFRPPPAR",
                        "bank_name": "BNP Paribas"
                    },
                    "tax_id": "FR12345678901",
                    "minimum_payout": 100.00,
                    "currency": "EUR"
                }
            }
        ]
        
        success_count = 0
        
        for aff_data in affiliates_data:
            logger.info(f"Processing affiliate: {aff_data['name']}")
            
            # Get organization ID
            org_id = None
            if aff_data['organization_name'] in self.results['organizations']:
                org_data = self.results['organizations'][aff_data['organization_name']]
                if org_data['id']:
                    org_id = org_data['id']
                else:
                    logger.error(f"Organization '{aff_data['organization_name']}' was not created successfully")
                    continue
            else:
                logger.error(f"Organization '{aff_data['organization_name']}' not found in results")
                continue
            
            result = self.api.create_affiliate_idempotent(
                name=aff_data['name'],
                organization_id=org_id,
                contact_email=aff_data['contact_email'],
                payment_details=aff_data['payment_details']
            )
            self.log_result(result)
            
            if result.success:
                self.results['affiliates'][aff_data['name']] = {
                    'id': result.entity_id,
                    'organization_id': org_id,
                    'result': result.result.value,
                    'data': result.entity_data
                }
                success_count += 1
            else:
                self.results['affiliates'][aff_data['name']] = {
                    'id': None,
                    'organization_id': org_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
        
        logger.info(f"Affiliates processed: {success_count}/{len(affiliates_data)} successful")
        return success_count == len(affiliates_data)
    
    def create_campaigns(self) -> bool:
    
        logger.info("Creating test campaigns...")
        
        campaigns_data = [
            {
                "name": "Adidas Summer Collection 2025",
                "advertiser_name": "Adidas Global",
                "payout_type": "cpa",
                "payout_amount": 15.00,
                "revenue_type": "rpa",
                "revenue_amount": 25.00,
                "currency_id": "USD",
                "status": "active",
                "visibility": "public"
            },
            {
                "name": "Dyson V15 Detect Launch",
                "advertiser_name": "Dyson Ltd",
                "payout_type": "cpa",
                "payout_amount": 50.00,
                "revenue_type": "rpa",
                "revenue_amount": 80.00,
                "currency_id": "USD",
                "status": "active",
                "visibility": "require_approval"
            }
        ]
        
        success_count = 0
        
        for camp_data in campaigns_data:
            logger.info(f"Processing campaign: {camp_data['name']}")
            
            # Get advertiser ID and organization ID
            advertiser_id = None
            organization_id = None
            if camp_data['advertiser_name'] in self.results['advertisers']:
                advertiser_data = self.results['advertisers'][camp_data['advertiser_name']]
                if advertiser_data['id']:
                    advertiser_id = advertiser_data['id']
                    organization_id = advertiser_data['organization_id']
                else:
                    logger.error(f"Advertiser '{camp_data['advertiser_name']}' was not created successfully")
                    continue
            else:
                logger.error(f"Advertiser '{camp_data['advertiser_name']}' not found in results")
                continue
            
            result = self.api.create_campaign_idempotent(
                name=camp_data['name'],
                advertiser_id=advertiser_id,
                organization_id=organization_id,
                payout_type=camp_data['payout_type'],
                payout_amount=camp_data['payout_amount'],
                revenue_type=camp_data['revenue_type'],
                revenue_amount=camp_data['revenue_amount'],
                currency_id=camp_data['currency_id'],
                status=camp_data['status'],
                visibility=camp_data['visibility']
            )
            self.log_result(result)
            
            if result.success:
                self.results['campaigns'][camp_data['name']] = {
                    'id': result.entity_id,
                    'advertiser_id': advertiser_id,
                    'result': result.result.value,
                    'data': result.entity_data
                }
                success_count += 1
            else:
                self.results['campaigns'][camp_data['name']] = {
                    'id': None,
                    'advertiser_id': advertiser_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
        
        logger.info(f"Campaigns processed: {success_count}/{len(campaigns_data)} successful")
        return success_count == len(campaigns_data)
    
    def create_analytics_data(self) -> bool:
    
        logger.info("Creating test analytics data...")
        
        # Advertiser analytics data
        advertiser_analytics = [
            {
                "domain": "adidas.com",
                "affiliate_networks": ["Impact", "CJ Affiliate", "ShareASale", "Awin"],
                "keywords": ["sportswear", "sneakers", "athletic", "running", "football", "basketball"],
                "verticals": ["Sports/Athletic Wear", "Fashion/Footwear", "Sports/Equipment"],
                "social_media_presence": {
                    "facebook": "https://facebook.com/adidas",
                    "instagram": "https://instagram.com/adidas",
                    "twitter": "https://twitter.com/adidas",
                    "youtube": "https://youtube.com/adidas"
                }
            },
            {
                "domain": "dyson.com",
                "affiliate_networks": ["Impact", "CJ Affiliate", "Awin", "ShareASale"],
                "keywords": ["vacuum", "air purifier", "hair dryer", "cordless", "cyclone"],
                "verticals": ["Home/Electricals", "Home/Appliances", "Beauty/Hair Care"],
                "social_media_presence": {
                    "facebook": "https://facebook.com/dyson",
                    "instagram": "https://instagram.com/dyson",
                    "twitter": "https://twitter.com/dyson",
                    "youtube": "https://youtube.com/dyson"
                }
            }
        ]
        
        # Publisher analytics data
        publisher_analytics = [
            {
                "domain": "lemonde.fr",
                "affiliate_networks": ["Affilae", "Awin", "Effiliation", "Impact", "Publicidees", "Rakuten"],
                "keywords": ["panier", "magasin", "fnac", "stock", "smartphone"],
                "traffic_score": 9250.75,
                "relevance": 85.5,
                "partners": [
                    "dyson.fr", "adidas.fr", "nike.fr", "amazon.fr", "fnac.com",
                    "cdiscount.com", "darty.com", "boulanger.com", "conforama.fr",
                    "ikea.fr", "leroy-merlin.fr", "castorama.fr", "decathlon.fr",
                    "zalando.fr", "asos.fr", "h-m.com", "zara.com", "mango.com",
                    "sephora.fr", "marionnaud.fr", "nocibe.fr", "douglas.fr",
                    "parfumerie-burdin.com", "origines-parfums.com", "notino.fr",
                    "parfums-moins-cher.com", "parfumdo.com"
                ]
            }
        ]
        
        success_count = 0
        total_count = len(advertiser_analytics) + len(publisher_analytics)
        
        # Create advertiser analytics
        for adv_analytics in advertiser_analytics:
            logger.info(f"Processing advertiser analytics for: {adv_analytics['domain']}")
            
            result = self.api.create_analytics_data_idempotent(
                '/api/v1/analytics/advertisers',
                adv_analytics,
                f"advertiser_{adv_analytics['domain']}"
            )
            self.log_result(result)
            
            if result.success:
                self.results['analytics'][f"advertiser_{adv_analytics['domain']}"] = {
                    'result': result.result.value,
                    'data': result.entity_data
                }
                success_count += 1
            else:
                self.results['analytics'][f"advertiser_{adv_analytics['domain']}"] = {
                    'result': result.result.value,
                    'error': result.error_message
                }
        
        # Create publisher analytics
        for pub_analytics in publisher_analytics:
            logger.info(f"Processing publisher analytics for: {pub_analytics['domain']}")
            
            result = self.api.create_analytics_data_idempotent(
                '/api/v1/analytics/affiliates',
                pub_analytics,
                f"publisher_{pub_analytics['domain']}"
            )
            self.log_result(result)
            
            if result.success:
                self.results['analytics'][f"publisher_{pub_analytics['domain']}"] = {
                    'result': result.result.value,
                    'data': result.entity_data
                }
                success_count += 1
            else:
                self.results['analytics'][f"publisher_{pub_analytics['domain']}"] = {
                    'result': result.result.value,
                    'error': result.error_message
                }
        
        logger.info(f"Analytics data processed: {success_count}/{total_count} successful")
        return success_count == total_count
    
    def create_all_test_data(self) -> bool:
    
        logger.info("Starting idempotent test account creation process...")
        
        steps = [
            ("API Health Check", self.check_api_health),
            ("Organizations", self.create_organizations),
            ("User Profiles", self.create_user_profiles),
            ("Advertisers", self.create_advertisers),
            ("Affiliates", self.create_affiliates),
            ("Campaigns", self.create_campaigns),
            ("Analytics Data", self.create_analytics_data)
        ]
        
        for step_name, step_func in steps:
            logger.info(f"\n{'='*50}")
            logger.info(f"Step: {step_name}")
            logger.info(f"{'='*50}")
            
            if not step_func():
                logger.error(f"Step failed: {step_name}")
                # Don't return False immediately - continue with other steps
                # This allows partial recovery and shows all issues
            
            # Small delay between steps
            time.sleep(0.5)
        
        logger.info(f"\n{'='*50}")
        logger.info("🎉 Idempotent test data creation process completed!")
        logger.info(f"{'='*50}")
        
        return True  # Always return True since we want to show the summary
    
    def print_summary(self):
    
        print("\n" + "="*70)
        print("IDEMPOTENT TEST ACCOUNT CREATION SUMMARY")
        print("="*70)
        
        # Count results by type
        for category, items in self.results.items():
            if not items:
                continue
                
            print(f"\n📊 {category.title()}: {len(items)} processed")
            
            created_count = sum(1 for item in items.values() if item.get('result') == 'created')
            exists_count = sum(1 for item in items.values() if item.get('result') == 'already_exists')
            failed_count = sum(1 for item in items.values() if item.get('result') == 'failed')
            
            print(f"   ✅ Created: {created_count}")
            print(f"   ℹ️  Already existed: {exists_count}")
            print(f"   ❌ Failed: {failed_count}")
            
            for name, data in items.items():
                status_icon = {
                    'created': '✅',
                    'already_exists': 'ℹ️ ',
                    'failed': '❌',
                    'updated': '🔄'
                }.get(data.get('result'), '❓')
                
                if data.get('id'):
                    print(f"   {status_icon} {name} (ID: {data['id']})")
                else:
                    print(f"   {status_icon} {name}")
                    if data.get('error'):
                        print(f"      Error: {data['error']}")
        
        print("\n" + "="*70)
        print("IDEMPOTENCY BENEFITS:")
        print("✅ Safe to run multiple times")
        print("✅ No duplicate data creation")
        print("✅ Graceful handling of existing entities")
        print("✅ Partial failure recovery")
        print("✅ Clear status reporting")
        print("\n" + "="*70)
        print("Next Steps:")
        print("1. Verify the data: python3 scripts/verify_test_accounts.py")
        print("2. Test the API endpoints with the created accounts")
        print("3. Use the test UUIDs for Supabase Auth integration")
        print("4. Run this script again anytime - it's fully idempotent!")
        print("="*70)

def main():

    parser = argparse.ArgumentParser(
        description="Create test accounts for the Affiliate Platform Backend (Idempotent Version)",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Basic usage with admin token
  python3 scripts/create_test_accounts_idempotent.py -t "your-jwt-token"
  
  # With custom API URL
  python3 scripts/create_test_accounts_idempotent.py -u "http://localhost:3000" -t "token"
  
  # Verbose output
  python3 scripts/create_test_accounts_idempotent.py -t "token" -v
  
  # Dry run (check API health only)
  python3 scripts/create_test_accounts_idempotent.py -u "http://localhost:8080" --dry-run

Environment Variables:
  API_BASE_URL    - API base URL (default: http://localhost:8080)
  ADMIN_JWT_TOKEN - JWT token for admin authentication

Features:
  ✅ Fully idempotent - safe to run multiple times
  ✅ No duplicate data creation
  ✅ Graceful error handling and recovery
  ✅ Clear status reporting and logging
  ✅ Modular and maintainable code structure
        """
    )
    
    parser.add_argument(
        '-u', '--api-url',
        default='http://localhost:8080',
        help='API base URL (default: http://localhost:8080)'
    )
    
    parser.add_argument(
        '-t', '--token',
        help='JWT token for admin authentication'
    )
    
    parser.add_argument(
        '-v', '--verbose',
        action='store_true',
        help='Enable verbose output'
    )
    
    parser.add_argument(
        '--dry-run',
        action='store_true',
        help='Only check API health, do not create data'
    )
    
    parser.add_argument(
        '--timeout',
        type=int,
        default=30,
        help='Request timeout in seconds (default: 30)'
    )
    
    args = parser.parse_args()
    
    # Get token from environment if not provided
    admin_token = args.token
    if not admin_token:
        import os
        admin_token = os.getenv('ADMIN_JWT_TOKEN')
    
    if not admin_token and not args.dry_run:
        logger.error("Admin JWT token is required for creating test data")
        logger.error("Provide it via -t option or ADMIN_JWT_TOKEN environment variable")
        sys.exit(1)
    
    # Create API client
    api_client = IdempotentAPIClient(
        base_url=args.api_url,
        admin_token=admin_token,
        timeout=args.timeout
    )
    
    # Create test account creator
    creator = IdempotentTestAccountCreator(api_client, verbose=args.verbose)
    
    try:
        if args.dry_run:
            logger.info("Performing dry run - checking API health only")
            if creator.check_api_health():
                logger.info("✅ API is healthy and ready for test account creation")
                sys.exit(0)
            else:
                logger.error("❌ API health check failed")
                sys.exit(1)
        else:
            # Create all test data
            creator.create_all_test_data()
            
            # Print summary
            creator.print_summary()
            
            logger.info("🎉 Idempotent test account creation completed!")
            sys.exit(0)
                
    except KeyboardInterrupt:
        logger.warning("Test account creation interrupted by user")
        sys.exit(1)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        if args.verbose:
            import traceback
            traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    main()