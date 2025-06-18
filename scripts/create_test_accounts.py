#!/usr/bin/env python3
"""
Simplified Test Account Creation Script
Creates test data using only the specified organizations and user IDs.
"""

import json
import logging
import requests
import time
from dataclasses import dataclass
from enum import Enum
from typing import List, Optional, Dict, Any

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class OperationResult(Enum):
    SUCCESS = "success"
    ALREADY_EXISTS = "already_exists"
    ERROR = "error"

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
    result: OperationResult
    data: Optional[Dict[str, Any]] = None
    error_message: Optional[str] = None

class IdempotentAPIClient:
    def __init__(self, base_url: str, admin_token: str = None):
        self.base_url = base_url.rstrip('/') + '/api/v1'
        self.admin_token = admin_token
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        })
        
        if self.admin_token:
            self.session.headers.update({
                'Authorization': f'Bearer {self.admin_token}'
            })

    def create_organization_idempotent(self, name: str, org_type: str) -> IdempotentResult:
        """Create organization if it doesn't exist"""
        try:
            # Check if organization exists
            response = self.session.get(f"{self.base_url}/organizations")
            if response.status_code == 200:
                orgs = response.json()
                for org in orgs:
                    if org.get('name') == name:
                        return IdempotentResult(
                            result=OperationResult.ALREADY_EXISTS,
                            data=org
                        )
            
            # Create organization
            data = {"name": name, "type": org_type}
            response = self.session.post(f"{self.base_url}/organizations", json=data)
            
            if response.status_code == 201:
                return IdempotentResult(
                    result=OperationResult.SUCCESS,
                    data=response.json()
                )
            else:
                return IdempotentResult(
                    result=OperationResult.ERROR,
                    error_message=f"HTTP {response.status_code}: {response.text}"
                )
                
        except Exception as e:
            return IdempotentResult(
                result=OperationResult.ERROR,
                error_message=str(e)
            )

    def create_user_profile_idempotent(self, user: TestUser, organization_id: str) -> IdempotentResult:
        """Create user profile if it doesn't exist"""
        try:
            # For user profiles, we'll use upsert since there's no list endpoint
            # and we want to create the profile with the specific UUID
            
            # Create user profile using upsert
            data = {
                "id": user.uuid,
                "email": user.email,
                "first_name": user.first_name,
                "last_name": user.last_name,
                "role_id": user.role_id,
                "organization_id": organization_id
            }
            response = self.session.post(f"{self.base_url}/profiles/upsert", json=data)
            
            if response.status_code in [200, 201]:
                result_type = OperationResult.ALREADY_EXISTS if response.status_code == 200 else OperationResult.SUCCESS
                return IdempotentResult(
                    result=result_type,
                    data=response.json()
                )
            else:
                return IdempotentResult(
                    result=OperationResult.ERROR,
                    error_message=f"HTTP {response.status_code}: {response.text}"
                )
                
        except Exception as e:
            return IdempotentResult(
                result=OperationResult.ERROR,
                error_message=str(e)
            )

    def create_advertiser_idempotent(self, name: str, organization_id: str, contact_email: str, billing_details: dict) -> IdempotentResult:
        """Create advertiser if it doesn't exist"""
        try:
            # Check if advertiser exists by listing organization's advertisers
            response = self.session.get(f"{self.base_url}/organizations/{organization_id}/advertisers")
            if response.status_code == 200:
                advertisers = response.json()
                if advertisers:  # Check if advertisers is not None
                    for adv in advertisers:
                        if adv.get('name') == name:
                            return IdempotentResult(
                                result=OperationResult.ALREADY_EXISTS,
                                data=adv
                            )
            
            # Create advertiser
            data = {
                "name": name,
                "organization_id": int(organization_id),
                "contact_email": contact_email,
                "billing_details": billing_details
            }
            response = self.session.post(f"{self.base_url}/advertisers", json=data)
            
            if response.status_code == 201:
                return IdempotentResult(
                    result=OperationResult.SUCCESS,
                    data=response.json()
                )
            else:
                return IdempotentResult(
                    result=OperationResult.ERROR,
                    error_message=f"HTTP {response.status_code}: {response.text}"
                )
                
        except Exception as e:
            return IdempotentResult(
                result=OperationResult.ERROR,
                error_message=str(e)
            )

    def create_affiliate_idempotent(self, name: str, organization_id: str, contact_email: str, **kwargs) -> IdempotentResult:
        """Create affiliate if it doesn't exist"""
        try:
            # Check if affiliate exists by listing organization's affiliates
            response = self.session.get(f"{self.base_url}/organizations/{organization_id}/affiliates")
            if response.status_code == 200:
                affiliates = response.json()
                if affiliates:  # Check if affiliates is not None
                    for aff in affiliates:
                        if aff.get('name') == name:
                            return IdempotentResult(
                                result=OperationResult.ALREADY_EXISTS,
                                data=aff
                            )
            
            # Create affiliate
            data = {
                "name": name,
                "organization_id": int(organization_id),
                "contact_email": contact_email,
                **kwargs
            }
            response = self.session.post(f"{self.base_url}/affiliates", json=data)
            
            if response.status_code == 201:
                return IdempotentResult(
                    result=OperationResult.SUCCESS,
                    data=response.json()
                )
            else:
                return IdempotentResult(
                    result=OperationResult.ERROR,
                    error_message=f"HTTP {response.status_code}: {response.text}"
                )
                
        except Exception as e:
            return IdempotentResult(
                result=OperationResult.ERROR,
                error_message=str(e)
            )

    def create_campaign_idempotent(self, name: str, advertiser_id: str, organization_id: str, payout_amount: float, payout_currency: str, **kwargs) -> IdempotentResult:
        """Create campaign if it doesn't exist"""
        try:
            # Check if campaign exists by listing organization's campaigns
            response = self.session.get(f"{self.base_url}/organizations/{organization_id}/campaigns")
            if response.status_code == 200:
                campaigns_response = response.json()
                if campaigns_response and 'campaigns' in campaigns_response:
                    campaigns = campaigns_response['campaigns']
                    for camp in campaigns:
                        if camp.get('name') == name:
                            return IdempotentResult(
                                result=OperationResult.ALREADY_EXISTS,
                                data=camp
                            )
            
            # Create campaign
            data = {
                "name": name,
                "organization_id": int(organization_id),
                "advertiser_id": int(advertiser_id),
                "payout_amount": payout_amount,
                "payout_currency": payout_currency,
                **kwargs
            }
            response = self.session.post(f"{self.base_url}/campaigns", json=data)
            
            if response.status_code == 201:
                return IdempotentResult(
                    result=OperationResult.SUCCESS,
                    data=response.json()
                )
            else:
                return IdempotentResult(
                    result=OperationResult.ERROR,
                    error_message=f"HTTP {response.status_code}: {response.text}"
                )
                
        except Exception as e:
            return IdempotentResult(
                result=OperationResult.ERROR,
                error_message=str(e)
            )

class SimpleTestAccountCreator:
    def __init__(self, api_client: IdempotentAPIClient, verbose: bool = False):
        self.api = api_client
        self.verbose = verbose
        self.results = {
            'organizations': {},
            'profiles': {},
            'advertisers': {},
            'affiliates': {},
            'campaigns': {}
        }
        
        if verbose:
            logger.setLevel(logging.DEBUG)

    def get_test_organizations(self) -> List[TestOrganization]:
        """Create test organizations using only the specified entities"""
        return [
            # Platform Owner (already exists in seed data)
            TestOrganization(
                name="rolinko",
                type="platform_owner",
                description="Platform administration organization"
            ),
            
            # Advertisers
            TestOrganization(
                name="Adidas",
                type="advertiser",
                description="Global sportswear brand"
            ),
            TestOrganization(
                name="Dyson",
                type="advertiser",
                description="British technology company - home appliances"
            ),
            
            # Affiliates
            TestOrganization(
                name="Le Monde",
                type="affiliate",
                description="French daily newspaper"
            )
        ]

    def get_test_users(self) -> List[TestUser]:
        """Create test users using only the specified user IDs and organizations"""
        return [
            # Skip Platform Admin (already exists in seed data)
            # TestUser(
            #     uuid="4cbe2452-88aa-4429-9145-b527d9eebfbf",
            #     email="admin@rolinko.com",
            #     first_name="Platform",
            #     last_name="Administrator",
            #     role_name="Admin",
            #     role_id=1,
            #     organization_name="rolinko"
            # ),
            
            # Advertiser Managers
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
            
            # Affiliate Managers
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
        """Create test organizations"""
        logger.info("Creating test organizations...")
        
        organizations = self.get_test_organizations()
        success_count = 0
        
        for org in organizations:
            logger.info(f"Processing organization: {org.name} ({org.type})")
            
            result = self.api.create_organization_idempotent(org.name, org.type)
            
            if result.result in [OperationResult.SUCCESS, OperationResult.ALREADY_EXISTS]:
                self.results['organizations'][org.name] = {
                    'id': result.data['organization_id'],
                    'result': result.result.value
                }
                success_count += 1
                logger.info(f"✅ Organization '{org.name}': {result.result.value}")
            else:
                self.results['organizations'][org.name] = {
                    'id': None,
                    'result': result.result.value,
                    'error': result.error_message
                }
                logger.error(f"❌ Organization '{org.name}': {result.error_message}")
        
        logger.info(f"Organizations processed: {success_count}/{len(organizations)} successful")
        return success_count == len(organizations)

    def create_user_profiles(self) -> bool:
        """Create test user profiles"""
        logger.info("Creating test user profiles...")
        
        users = self.get_test_users()
        success_count = 0
        
        for user in users:
            logger.info(f"Processing profile for: {user.email}")
            
            # Get organization ID
            org_id = None
            if user.organization_name in self.results['organizations']:
                org_id = self.results['organizations'][user.organization_name]['id']
            
            if not org_id:
                logger.error(f"❌ Organization '{user.organization_name}' not found for user {user.email}")
                continue
            
            result = self.api.create_user_profile_idempotent(user, org_id)
            
            if result.result in [OperationResult.SUCCESS, OperationResult.ALREADY_EXISTS]:
                self.results['profiles'][user.email] = {
                    'id': result.data['id'],
                    'organization_id': org_id,
                    'result': result.result.value
                }
                success_count += 1
                logger.info(f"✅ Profile '{user.email}': {result.result.value}")
            else:
                self.results['profiles'][user.email] = {
                    'id': None,
                    'organization_id': org_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
                logger.error(f"❌ Profile '{user.email}': {result.error_message}")
        
        logger.info(f"Profiles processed: {success_count}/{len(users)} successful")
        return success_count == len(users)

    def create_advertisers(self) -> bool:
        """Create test advertisers"""
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
                org_id = self.results['organizations'][adv_data['organization_name']]['id']
            
            if not org_id:
                logger.error(f"❌ Organization '{adv_data['organization_name']}' not found")
                continue
            
            result = self.api.create_advertiser_idempotent(
                name=adv_data['name'],
                organization_id=org_id,
                contact_email=adv_data['contact_email'],
                billing_details=adv_data['billing_details']
            )
            
            if result.result in [OperationResult.SUCCESS, OperationResult.ALREADY_EXISTS]:
                self.results['advertisers'][adv_data['name']] = {
                    'id': result.data['advertiser_id'],
                    'organization_id': org_id,
                    'result': result.result.value
                }
                success_count += 1
                logger.info(f"✅ Advertiser '{adv_data['name']}': {result.result.value}")
            else:
                self.results['advertisers'][adv_data['name']] = {
                    'id': None,
                    'organization_id': org_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
                logger.error(f"❌ Advertiser '{adv_data['name']}': {result.error_message}")
        
        logger.info(f"Advertisers processed: {success_count}/{len(advertisers_data)} successful")
        return success_count == len(advertisers_data)

    def create_affiliates(self) -> bool:
        """Create test affiliates"""
        logger.info("Creating test affiliates...")
        
        affiliates_data = [
            {
                "name": "Le Monde Digital",
                "organization_name": "Le Monde",
                "contact_email": "rolinko@lemonde.fr",
                "payment_details": {
                    "bank_account": "FR1420041010050500013M02606",
                    "routing_number": "PSSTFRPPPAR",
                    "payment_method": "bank_transfer"
                },
                "status": "active",
                "internal_notes": "Premium French news publisher",
                "default_currency_id": "EUR",
                "contact_address": {
                    "address1": "67-69 Avenue Pierre Mendès France",
                    "city": "Paris",
                    "region_code": "IDF",
                    "country_code": "FR",
                    "zip_postal_code": "75013"
                },
                "billing_info": {
                    "company_name": "Le Monde SA",
                    "tax_id": "FR12345678901"
                },
                "labels": ["news", "premium", "french"],
                "invoice_amount_threshold": 500.00,
                "default_payment_terms": 30
            }
        ]
        
        success_count = 0
        
        for aff_data in affiliates_data:
            logger.info(f"Processing affiliate: {aff_data['name']}")
            
            # Get organization ID
            org_id = None
            if aff_data['organization_name'] in self.results['organizations']:
                org_id = self.results['organizations'][aff_data['organization_name']]['id']
            
            if not org_id:
                logger.error(f"❌ Organization '{aff_data['organization_name']}' not found")
                continue
            
            # Extract affiliate-specific fields
            affiliate_kwargs = {k: v for k, v in aff_data.items() 
                              if k not in ['name', 'organization_name', 'contact_email']}
            
            result = self.api.create_affiliate_idempotent(
                name=aff_data['name'],
                organization_id=org_id,
                contact_email=aff_data['contact_email'],
                **affiliate_kwargs
            )
            
            if result.result in [OperationResult.SUCCESS, OperationResult.ALREADY_EXISTS]:
                self.results['affiliates'][aff_data['name']] = {
                    'id': result.data['affiliate_id'],
                    'organization_id': org_id,
                    'result': result.result.value
                }
                success_count += 1
                logger.info(f"✅ Affiliate '{aff_data['name']}': {result.result.value}")
            else:
                self.results['affiliates'][aff_data['name']] = {
                    'id': None,
                    'organization_id': org_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
                logger.error(f"❌ Affiliate '{aff_data['name']}': {result.error_message}")
        
        logger.info(f"Affiliates processed: {success_count}/{len(affiliates_data)} successful")
        return success_count == len(affiliates_data)

    def create_campaigns(self) -> bool:
        """Create test campaigns"""
        logger.info("Creating test campaigns...")
        
        campaigns_data = [
            {
                "name": "Adidas Holiday Collection 2025",
                "advertiser_name": "Adidas Global",
                "payout_amount": 15.00,
                "payout_currency": "EUR",
                "status": "active",
                "description": "Holiday season sportswear and lifestyle collection",
                "start_date": "2025-11-01T00:00:00Z",
                "end_date": "2025-12-31T23:59:59Z",
                "destination_url": "https://adidas.com/holiday-collection",
                "thumbnail_url": "https://assets.adidas.com/holiday-thumb.jpg"
            },
            {
                "name": "Dyson Christmas Gift Guide",
                "advertiser_name": "Dyson Ltd",
                "payout_amount": 25.00,
                "payout_currency": "GBP",
                "status": "active",
                "description": "Premium home appliances for holiday gifting",
                "start_date": "2025-11-15T00:00:00Z",
                "end_date": "2025-12-25T23:59:59Z",
                "destination_url": "https://dyson.com/christmas-gifts",
                "thumbnail_url": "https://assets.dyson.com/christmas-thumb.jpg"
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
                advertiser_id = advertiser_data['id']
                organization_id = advertiser_data['organization_id']
            
            if not advertiser_id or not organization_id:
                logger.error(f"❌ Advertiser '{camp_data['advertiser_name']}' not found")
                continue
            
            # Extract campaign-specific fields
            campaign_kwargs = {k: v for k, v in camp_data.items() 
                             if k not in ['name', 'advertiser_name', 'payout_amount', 'payout_currency']}
            
            result = self.api.create_campaign_idempotent(
                name=camp_data['name'],
                advertiser_id=advertiser_id,
                organization_id=organization_id,
                payout_amount=camp_data['payout_amount'],
                payout_currency=camp_data['payout_currency'],
                **campaign_kwargs
            )
            
            if result.result in [OperationResult.SUCCESS, OperationResult.ALREADY_EXISTS]:
                self.results['campaigns'][camp_data['name']] = {
                    'id': result.data['campaign_id'],
                    'advertiser_id': advertiser_id,
                    'result': result.result.value
                }
                success_count += 1
                logger.info(f"✅ Campaign '{camp_data['name']}': {result.result.value}")
            else:
                self.results['campaigns'][camp_data['name']] = {
                    'id': None,
                    'advertiser_id': advertiser_id,
                    'result': result.result.value,
                    'error': result.error_message
                }
                logger.error(f"❌ Campaign '{camp_data['name']}': {result.error_message}")
        
        logger.info(f"Campaigns processed: {success_count}/{len(campaigns_data)} successful")
        return success_count == len(campaigns_data)

    def run_all(self) -> bool:
        """Run all test data creation steps"""
        logger.info("🚀 Starting simplified test account creation...")
        
        steps = [
            ("Organizations", self.create_organizations),
            ("User Profiles", self.create_user_profiles),
            ("Advertisers", self.create_advertisers),
            ("Affiliates", self.create_affiliates),
            ("Campaigns", self.create_campaigns)
        ]
        
        all_success = True
        
        for step_name, step_func in steps:
            logger.info(f"\n📋 Step: {step_name}")
            success = step_func()
            if not success:
                logger.error(f"❌ Step '{step_name}' failed")
                all_success = False
            else:
                logger.info(f"✅ Step '{step_name}' completed successfully")
        
        # Print summary
        logger.info("\n📊 CREATION SUMMARY:")
        for category, items in self.results.items():
            successful = sum(1 for item in items.values() if item.get('result') in ['success', 'already_exists'])
            total = len(items)
            logger.info(f"  {category.title()}: {successful}/{total} successful")
        
        if all_success:
            logger.info("🎉 All test data created successfully!")
        else:
            logger.error("❌ Some steps failed. Check logs for details.")
        
        return all_success

def main():
    """Main function"""
    import sys
    import os
    import argparse
    
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Create simplified test accounts for affiliate platform')
    parser.add_argument(
        '--api-url',
        default='http://localhost:8080',
        help='Base URL for the API (default: http://localhost:8080)'
    )
    parser.add_argument(
        '-t', '--token',
        help='JWT token for admin authentication'
    )
    parser.add_argument(
        '--verbose',
        action='store_true',
        help='Enable verbose logging'
    )
    
    args = parser.parse_args()
    
    # Get token from environment if not provided
    admin_token = args.token
    if not admin_token:
        admin_token = os.getenv('ADMIN_JWT_TOKEN')
    
    if not admin_token:
        logger.error("Admin JWT token is required for creating test data")
        logger.error("Provide it via -t option or ADMIN_JWT_TOKEN environment variable")
        logger.error("\nTo generate a test token, run:")
        logger.error("  python scripts/generate_test_jwt.py")
        sys.exit(1)
    
    # Create API client
    api_client = IdempotentAPIClient(args.api_url, admin_token)
    
    # Create test account creator
    creator = SimpleTestAccountCreator(api_client, verbose=args.verbose)
    
    # Run all creation steps
    success = creator.run_all()
    
    # Exit with appropriate code
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()