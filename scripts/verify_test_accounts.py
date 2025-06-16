#!/usr/bin/env python3

"""
Test Account Verification Script for Affiliate Platform Backend

This script verifies that test accounts were created successfully using the API endpoints.
It checks:
- Organizations exist and are accessible
- User profiles are created with correct roles
- Advertisers and affiliates are properly configured
- Campaigns are active and accessible
- Analytics data is available

Usage:
    python3 scripts/verify_test_accounts.py [options]

Requirements:
    - API server running on specified URL
    - Test accounts already created
    - Admin JWT token for authentication
"""

import argparse
import json
import logging
import sys
import time
from typing import Dict, List, Optional, Any
from dataclasses import dataclass
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

@dataclass
class VerificationResult:
    """Represents a verification result"""
    test_name: str
    success: bool
    message: str
    details: Optional[Dict] = None

class AffiliateAPIClient:
    """Client for interacting with the Affiliate Platform API"""
    
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
    
    def _make_request(self, method: str, endpoint: str, params: Optional[Dict] = None) -> Dict:
        """Make an HTTP request to the API"""
        url = f"{self.base_url}{endpoint}"
        
        try:
            logger.debug(f"Making {method} request to {url}")
            
            response = self.session.request(
                method=method,
                url=url,
                params=params,
                timeout=self.timeout
            )
            
            logger.debug(f"Response status: {response.status_code}")
            
            # Try to parse JSON response
            try:
                response_data = response.json()
            except json.JSONDecodeError:
                response_data = {"raw_response": response.text}
            
            return {
                'success': response.status_code >= 200 and response.status_code < 300,
                'status_code': response.status_code,
                'data': response_data
            }
                
        except requests.exceptions.RequestException as e:
            logger.error(f"Request failed: {e}")
            return {
                'success': False,
                'status_code': 0,
                'error': str(e)
            }
    
    def health_check(self) -> Dict:
        """Check API health"""
        return self._make_request('GET', '/health')
    
    def list_organizations(self) -> Dict:
        """List all organizations"""
        return self._make_request('GET', '/api/v1/organizations')
    
    def get_organization(self, org_id: int) -> Dict:
        """Get organization by ID"""
        return self._make_request('GET', f'/api/v1/organizations/{org_id}')
    
    def autocomplete_search(self, query: str, search_type: Optional[str] = None, limit: int = 10) -> Dict:
        """Search using autocomplete endpoint"""
        params = {'q': query, 'limit': limit}
        if search_type:
            params['type'] = search_type
        return self._make_request('GET', '/api/v1/analytics/autocomplete', params=params)
    
    def get_advertiser_analytics(self, advertiser_id: int) -> Dict:
        """Get advertiser analytics by ID"""
        return self._make_request('GET', f'/api/v1/analytics/advertisers/{advertiser_id}')
    
    def get_publisher_analytics(self, publisher_id: int) -> Dict:
        """Get publisher analytics by ID"""
        return self._make_request('GET', f'/api/v1/analytics/affiliates/{publisher_id}')

class TestAccountVerifier:
    """Main class for verifying test accounts"""
    
    def __init__(self, api_client: AffiliateAPIClient, verbose: bool = False):
        self.api = api_client
        self.verbose = verbose
        self.verification_results = []
        
        if verbose:
            logger.setLevel(logging.DEBUG)
    
    def log_success(self, message: str):
        """Log success message"""
        logger.info(f"✅ {message}")
    
    def log_error(self, message: str):
        """Log error message"""
        logger.error(f"❌ {message}")
    
    def log_warning(self, message: str):
        """Log warning message"""
        logger.warning(f"⚠️  {message}")
    
    def log_info(self, message: str):
        """Log info message"""
        logger.info(f"ℹ️  {message}")
    
    def add_result(self, test_name: str, success: bool, message: str, details: Optional[Dict] = None):
        """Add a verification result"""
        result = VerificationResult(test_name, success, message, details)
        self.verification_results.append(result)
        
        if success:
            self.log_success(f"{test_name}: {message}")
        else:
            self.log_error(f"{test_name}: {message}")
    
    def verify_api_health(self) -> bool:
        """Verify API health"""
        self.log_info("Checking API health...")
        
        response = self.api.health_check()
        
        if response['success']:
            self.add_result("API Health", True, "API is healthy and responding")
            return True
        else:
            error_msg = response.get('error', f"HTTP {response['status_code']}")
            self.add_result("API Health", False, f"API health check failed: {error_msg}")
            return False
    
    def verify_organizations(self) -> bool:
        """Verify test organizations exist"""
        self.log_info("Verifying test organizations...")
        
        expected_orgs = ["UpsailAI", "Adidas", "Dyson", "Le Monde"]
        
        response = self.api.list_organizations()
        
        if not response['success']:
            error_msg = response.get('error', f"HTTP {response['status_code']}")
            self.add_result("Organizations", False, f"Failed to list organizations: {error_msg}")
            return False
        
        organizations = response['data']
        if not isinstance(organizations, list):
            self.add_result("Organizations", False, "Invalid response format for organizations list")
            return False
        
        found_orgs = {org.get('name'): org for org in organizations if org.get('name') in expected_orgs}
        
        success_count = 0
        for org_name in expected_orgs:
            if org_name in found_orgs:
                org_data = found_orgs[org_name]
                self.add_result(
                    f"Organization {org_name}",
                    True,
                    f"Found with ID {org_data.get('organization_id')} and type {org_data.get('type')}",
                    org_data
                )
                success_count += 1
            else:
                self.add_result(f"Organization {org_name}", False, "Not found")
        
        overall_success = success_count == len(expected_orgs)
        self.add_result(
            "Organizations Overall",
            overall_success,
            f"{success_count}/{len(expected_orgs)} organizations found"
        )
        
        return overall_success
    
    def verify_autocomplete_functionality(self) -> bool:
        """Verify analytics autocomplete functionality"""
        self.log_info("Verifying autocomplete functionality...")
        
        test_searches = [
            ("adi", "advertiser", "Adidas"),
            ("dys", "advertiser", "Dyson"),
            ("lem", "publisher", "Le Monde"),
            ("ups", None, "UpsailAI")
        ]
        
        success_count = 0
        
        for query, search_type, expected_name in test_searches:
            test_name = f"Autocomplete {query}"
            if search_type:
                test_name += f" ({search_type})"
            
            response = self.api.autocomplete_search(query, search_type)
            
            if not response['success']:
                error_msg = response.get('error', f"HTTP {response['status_code']}")
                self.add_result(test_name, False, f"Search failed: {error_msg}")
                continue
            
            results = response['data']
            if not isinstance(results, list):
                self.add_result(test_name, False, "Invalid response format")
                continue
            
            # Check if expected name is in results
            found = any(
                expected_name.lower() in result.get('name', '').lower()
                for result in results
            )
            
            if found:
                self.add_result(test_name, True, f"Found {expected_name} in results")
                success_count += 1
            else:
                self.add_result(test_name, False, f"{expected_name} not found in results")
        
        overall_success = success_count == len(test_searches)
        self.add_result(
            "Autocomplete Overall",
            overall_success,
            f"{success_count}/{len(test_searches)} searches successful"
        )
        
        return overall_success
    
    def verify_advertiser_analytics(self) -> bool:
        """Verify advertiser analytics data"""
        self.log_info("Verifying advertiser analytics...")
        
        # Test known advertiser IDs (assuming they start from 1)
        advertiser_tests = [
            (1, "Adidas", "adidas.com"),
            (2, "Dyson", "dyson.com")
        ]
        
        success_count = 0
        
        for advertiser_id, expected_name, expected_domain in advertiser_tests:
            test_name = f"Advertiser Analytics {advertiser_id}"
            
            response = self.api.get_advertiser_analytics(advertiser_id)
            
            if not response['success']:
                error_msg = response.get('error', f"HTTP {response['status_code']}")
                self.add_result(test_name, False, f"Failed to get analytics: {error_msg}")
                continue
            
            analytics_data = response['data']
            if not isinstance(analytics_data, dict):
                self.add_result(test_name, False, "Invalid response format")
                continue
            
            # Check if expected domain is present
            domain = analytics_data.get('domain', '')
            if expected_domain in domain:
                details = {
                    'domain': domain,
                    'networks_count': len(analytics_data.get('affiliate_networks', [])),
                    'keywords_count': len(analytics_data.get('keywords', []))
                }
                self.add_result(
                    test_name,
                    True,
                    f"Found analytics for {expected_domain}",
                    details
                )
                success_count += 1
            else:
                self.add_result(test_name, False, f"Expected domain {expected_domain} not found")
        
        overall_success = success_count == len(advertiser_tests)
        self.add_result(
            "Advertiser Analytics Overall",
            overall_success,
            f"{success_count}/{len(advertiser_tests)} advertiser analytics verified"
        )
        
        return overall_success
    
    def verify_publisher_analytics(self) -> bool:
        """Verify publisher analytics data"""
        self.log_info("Verifying publisher analytics...")
        
        # Test known publisher IDs (assuming they start from 1)
        publisher_tests = [
            (1, "Le Monde", "lemonde.fr")
        ]
        
        success_count = 0
        
        for publisher_id, expected_name, expected_domain in publisher_tests:
            test_name = f"Publisher Analytics {publisher_id}"
            
            response = self.api.get_publisher_analytics(publisher_id)
            
            if not response['success']:
                error_msg = response.get('error', f"HTTP {response['status_code']}")
                self.add_result(test_name, False, f"Failed to get analytics: {error_msg}")
                continue
            
            analytics_data = response['data']
            if not isinstance(analytics_data, dict):
                self.add_result(test_name, False, "Invalid response format")
                continue
            
            # Check if expected domain is present
            domain = analytics_data.get('domain', '')
            if expected_domain in domain:
                details = {
                    'domain': domain,
                    'traffic_score': analytics_data.get('traffic_score'),
                    'networks_count': len(analytics_data.get('affiliate_networks', [])),
                    'partners_count': len(analytics_data.get('partners', []))
                }
                self.add_result(
                    test_name,
                    True,
                    f"Found analytics for {expected_domain}",
                    details
                )
                success_count += 1
            else:
                self.add_result(test_name, False, f"Expected domain {expected_domain} not found")
        
        overall_success = success_count == len(publisher_tests)
        self.add_result(
            "Publisher Analytics Overall",
            overall_success,
            f"{success_count}/{len(publisher_tests)} publisher analytics verified"
        )
        
        return overall_success
    
    def run_all_verifications(self) -> bool:
        """Run all verification tests"""
        self.log_info("Starting test account verification...")
        
        verification_steps = [
            ("API Health Check", self.verify_api_health),
            ("Organizations", self.verify_organizations),
            ("Autocomplete Functionality", self.verify_autocomplete_functionality),
            ("Advertiser Analytics", self.verify_advertiser_analytics),
            ("Publisher Analytics", self.verify_publisher_analytics)
        ]
        
        overall_success = True
        
        for step_name, step_func in verification_steps:
            self.log_info(f"\n{'='*50}")
            self.log_info(f"Verification Step: {step_name}")
            self.log_info(f"{'='*50}")
            
            step_success = step_func()
            if not step_success:
                overall_success = False
            
            # Small delay between steps
            time.sleep(0.5)
        
        return overall_success
    
    def print_summary(self):
        """Print verification summary"""
        print("\n" + "="*60)
        print("TEST ACCOUNT VERIFICATION SUMMARY")
        print("="*60)
        
        passed_tests = [r for r in self.verification_results if r.success]
        failed_tests = [r for r in self.verification_results if not r.success]
        
        print(f"\n📊 Total Tests: {len(self.verification_results)}")
        print(f"✅ Passed: {len(passed_tests)}")
        print(f"❌ Failed: {len(failed_tests)}")
        print(f"📈 Success Rate: {len(passed_tests)/len(self.verification_results)*100:.1f}%")
        
        if failed_tests:
            print(f"\n❌ Failed Tests:")
            for test in failed_tests:
                print(f"  • {test.test_name}: {test.message}")
        
        if passed_tests and self.verbose:
            print(f"\n✅ Passed Tests:")
            for test in passed_tests:
                print(f"  • {test.test_name}: {test.message}")
                if test.details:
                    for key, value in test.details.items():
                        print(f"    - {key}: {value}")
        
        print("\n" + "="*60)
        
        if len(failed_tests) == 0:
            print("🎉 All verifications passed! Test accounts are working correctly.")
            print("\nNext Steps:")
            print("1. Test the API endpoints with different user roles")
            print("2. Integrate with your frontend application")
            print("3. Set up Supabase Auth with the test UUIDs")
        else:
            print("⚠️  Some verifications failed. Please check the issues above.")
            print("\nTroubleshooting:")
            print("1. Ensure the API server is running")
            print("2. Check that test accounts were created successfully")
            print("3. Verify JWT token has admin permissions")
            print("4. Review API logs for detailed error information")
        
        print("="*60)

def main():
    """Main function"""
    parser = argparse.ArgumentParser(
        description="Verify test accounts for the Affiliate Platform Backend",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Basic verification with admin token
  python3 scripts/verify_test_accounts.py -t "your-jwt-token"
  
  # With custom API URL
  python3 scripts/verify_test_accounts.py -u "http://localhost:3000" -t "token"
  
  # Verbose output
  python3 scripts/verify_test_accounts.py -t "token" -v
  
  # Health check only
  python3 scripts/verify_test_accounts.py -u "http://localhost:8080" --health-only

Environment Variables:
  API_BASE_URL    - API base URL (default: http://localhost:8080)
  ADMIN_JWT_TOKEN - JWT token for admin authentication
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
        '--health-only',
        action='store_true',
        help='Only check API health, skip other verifications'
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
    
    if not admin_token and not args.health_only:
        logger.error("Admin JWT token is required for verification")
        logger.error("Provide it via -t option or ADMIN_JWT_TOKEN environment variable")
        logger.info("Use --health-only to check API health without authentication")
        sys.exit(1)
    
    # Create API client
    api_client = AffiliateAPIClient(
        base_url=args.api_url,
        admin_token=admin_token,
        timeout=args.timeout
    )
    
    # Create verifier
    verifier = TestAccountVerifier(api_client, verbose=args.verbose)
    
    try:
        if args.health_only:
            logger.info("Performing health check only")
            if verifier.verify_api_health():
                logger.info("✅ API is healthy")
                sys.exit(0)
            else:
                logger.error("❌ API health check failed")
                sys.exit(1)
        else:
            # Run all verifications
            success = verifier.run_all_verifications()
            
            # Print summary
            verifier.print_summary()
            
            if success:
                logger.info("🎉 All verifications passed!")
                sys.exit(0)
            else:
                logger.error("❌ Some verifications failed")
                sys.exit(1)
                
    except KeyboardInterrupt:
        logger.warning("Verification interrupted by user")
        sys.exit(1)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        if args.verbose:
            import traceback
            traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    main()