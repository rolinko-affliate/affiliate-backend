#!/usr/bin/env python3
"""
Organization-Specific Mock Data Generator

This script takes a JWT token as input, queries the user's profile and organization,
then generates hundreds of click and conversion events for that specific organization.
The generated data is saved to organization-specific folders.

Usage:
    python generate_org_mock_data.py --jwt-token <token> [--base-url <url>] [--output-dir <dir>] [--num-clicks <num>]
"""

import argparse
import csv
import json
import logging
import os
import random
import requests
import uuid
from datetime import datetime, timedelta
from dataclasses import dataclass, asdict
from typing import List, Dict, Optional, Any
from faker import Faker

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Initialize Faker for generating realistic data
fake = Faker()

@dataclass
class UserProfile:
    id: str
    organization_id: Optional[int]
    role_id: int
    role_name: str
    email: str
    first_name: Optional[str] = None
    last_name: Optional[str] = None
    created_at: Optional[str] = None
    updated_at: Optional[str] = None

@dataclass
class Organization:
    organization_id: int
    name: str
    type: str
    created_at: str
    updated_at: str

@dataclass
class Campaign:
    campaign_id: int
    organization_id: int
    advertiser_id: int
    name: str
    status: str
    visibility: Optional[str] = None
    currency_id: Optional[str] = None
    created_at: Optional[str] = None
    updated_at: Optional[str] = None
    description: Optional[str] = None
    destination_url: Optional[str] = None

@dataclass
class Advertiser:
    advertiser_id: int
    organization_id: int
    name: str
    status: str
    created_at: Optional[str] = None
    updated_at: Optional[str] = None
    contact_email: Optional[str] = None

@dataclass
class Affiliate:
    affiliate_id: int
    organization_id: int
    name: str
    status: str
    created_at: Optional[str] = None
    updated_at: Optional[str] = None
    contact_email: Optional[str] = None

@dataclass
class ClickEvent:
    organization_id: int
    id: str
    timestamp: str
    campaign_id: int
    campaign_name: str
    offer_id: str
    offer_name: str
    affiliate_id: int
    affiliate_name: str
    ip_address: str
    user_agent: str
    country: str
    region: str
    city: str
    referrer_url: str
    landing_page_url: str
    sub1: str
    sub2: str
    sub3: str
    converted: bool
    conversion_id: Optional[str] = None

@dataclass
class ConversionEvent:
    organization_id: int
    id: str
    timestamp: str
    transaction_id: str
    campaign_id: int
    campaign_name: str
    offer_id: str
    offer_name: str
    status: str
    payout: float
    currency: str
    affiliate_id: int
    affiliate_name: str
    click_id: str
    conversion_value: float
    sub1: str
    sub2: str
    sub3: str

@dataclass
class PerformanceSummary:
    organization_id: int
    total_clicks: int
    total_conversions: int
    total_revenue: float
    conversion_rate: float
    average_revenue: float
    click_through_rate: float
    total_impressions: int

@dataclass
class DailyPerformanceReport:
    organization_id: int
    date: str
    campaign_id: int
    campaign_name: str
    clicks: int
    impressions: int
    conversions: int
    revenue: float
    conversion_rate: float
    click_through_rate: float
    payouts: float

@dataclass
class CampaignPerformance:
    campaign_id: int
    organization_id: int
    name: str
    clicks: int
    conversions: int
    revenue: float
    conversion_rate: float
    status: str
    priority: int
    tier: int

class APIClient:
    def __init__(self, base_url: str, jwt_token: str):
        self.base_url = base_url.rstrip('/') + '/api/v1'
        self.jwt_token = jwt_token
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': f'Bearer {jwt_token}'
        })

    def get_user_profile(self) -> UserProfile:
        """Get the current user's profile"""
        response = self.session.get(f"{self.base_url}/users/me")
        response.raise_for_status()
        data = response.json()
        
        # Extract only the fields we need for UserProfile
        profile_data = {
            'id': data.get('id'),
            'organization_id': data.get('organization_id'),
            'role_id': data.get('role_id'),
            'role_name': data.get('role_name'),
            'email': data.get('email'),
            'first_name': data.get('first_name'),
            'last_name': data.get('last_name'),
            'created_at': data.get('created_at'),
            'updated_at': data.get('updated_at')
        }
        
        return UserProfile(**profile_data)

    def get_organization(self, org_id: int) -> Organization:
        """Get organization details"""
        response = self.session.get(f"{self.base_url}/organizations/{org_id}")
        response.raise_for_status()
        data = response.json()
        
        # Filter out fields that don't exist in the Organization dataclass
        import inspect
        sig = inspect.signature(Organization)
        filtered_data = {k: v for k, v in data.items() if k in sig.parameters}
        return Organization(**filtered_data)

    def get_organization_campaigns(self, org_id: int) -> List[Campaign]:
        """Get campaigns for an organization"""
        response = self.session.get(f"{self.base_url}/organizations/{org_id}/campaigns")
        response.raise_for_status()
        data = response.json()
        
        campaigns = []
        if 'campaigns' in data:
            for campaign_data in data['campaigns']:
                # Filter out fields that don't exist in the Campaign dataclass
                import inspect
                sig = inspect.signature(Campaign)
                filtered_data = {k: v for k, v in campaign_data.items() if k in sig.parameters}
                campaigns.append(Campaign(**filtered_data))
        return campaigns

    def get_organization_advertisers(self, org_id: int) -> List[Advertiser]:
        """Get advertisers for an organization"""
        try:
            response = self.session.get(f"{self.base_url}/organizations/{org_id}/advertisers")
            response.raise_for_status()
            data = response.json()
            
            advertisers = []
            if data:
                for advertiser_data in data:
                    # Filter out fields that don't exist in the Advertiser dataclass
                    import inspect
                    sig = inspect.signature(Advertiser)
                    filtered_data = {k: v for k, v in advertiser_data.items() if k in sig.parameters}
                    advertisers.append(Advertiser(**filtered_data))
            return advertisers
        except requests.exceptions.HTTPError as e:
            if e.response.status_code == 404:
                logger.warning(f"No advertisers found for organization {org_id}")
                return []
            raise

    def get_organization_affiliates(self, org_id: int) -> List[Affiliate]:
        """Get affiliates for an organization"""
        try:
            response = self.session.get(f"{self.base_url}/organizations/{org_id}/affiliates")
            response.raise_for_status()
            data = response.json()
            
            affiliates = []
            if data:
                for affiliate_data in data:
                    # Filter out fields that don't exist in the Affiliate dataclass
                    import inspect
                    sig = inspect.signature(Affiliate)
                    filtered_data = {k: v for k, v in affiliate_data.items() if k in sig.parameters}
                    affiliates.append(Affiliate(**filtered_data))
            return affiliates
        except requests.exceptions.HTTPError as e:
            if e.response.status_code == 404:
                logger.warning(f"No affiliates found for organization {org_id}")
                return []
            raise

class MockDataGenerator:
    def __init__(self, organization: Organization, campaigns: List[Campaign], 
                 advertisers: List[Advertiser], affiliates: List[Affiliate]):
        self.organization = organization
        self.campaigns = campaigns
        self.advertisers = advertisers
        self.affiliates = affiliates
        
        # Create default data if none exists
        if not self.campaigns:
            self.campaigns = self._create_default_campaigns()
        if not self.advertisers:
            self.advertisers = self._create_default_advertisers()
        if not self.affiliates:
            self.affiliates = self._create_default_affiliates()

        # Pre-generate some common data
        self.user_agents = [
            'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
            'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
            'Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1',
            'Mozilla/5.0 (Android 11; Mobile; rv:68.0) Gecko/68.0 Firefox/88.0',
            'Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1'
        ]
        
        self.referrer_sources = [
            'https://google.com/search?q=',
            'https://facebook.com/',
            'https://instagram.com/',
            'https://youtube.com/watch?v=',
            'https://tiktok.com/@',
            'https://pinterest.com/',
            'https://twitter.com/',
            'https://bing.com/search?q=',
            'https://email-newsletter.com',
            'https://affiliate-site.com'
        ]
        
        self.countries = ['US', 'CA', 'GB', 'DE', 'FR', 'AU', 'JP', 'BR', 'IN', 'MX']
        self.conversion_statuses = ['approved', 'pending', 'rejected']
        self.currencies = ['USD', 'EUR', 'GBP', 'CAD', 'AUD']

    def _create_default_campaigns(self) -> List[Campaign]:
        """Create default campaigns if none exist"""
        default_campaigns = []
        for i in range(3):
            campaign = Campaign(
                campaign_id=10000 + i,
                organization_id=self.organization.organization_id,
                advertiser_id=1000 + i,
                name=f"{self.organization.name} Campaign {i+1}",
                description=f"Default campaign {i+1} for {self.organization.name}",
                status="active"
            )
            default_campaigns.append(campaign)
        return default_campaigns

    def _create_default_advertisers(self) -> List[Advertiser]:
        """Create default advertisers if none exist"""
        default_advertisers = []
        for i in range(2):
            advertiser = Advertiser(
                advertiser_id=1000 + i,
                organization_id=self.organization.organization_id,
                name=f"{self.organization.name} Advertiser {i+1}",
                contact_email=f"advertiser{i+1}@{self.organization.name.lower().replace(' ', '')}.com"
            )
            default_advertisers.append(advertiser)
        return default_advertisers

    def _create_default_affiliates(self) -> List[Affiliate]:
        """Create default affiliates if none exist"""
        default_affiliates = []
        affiliate_names = ['TechReviews Pro', 'Beauty Influencer', 'Home & Garden', 'Lifestyle Blog', 'Deal Hunter']
        for i, name in enumerate(affiliate_names):
            affiliate = Affiliate(
                affiliate_id=2000 + i,
                organization_id=self.organization.organization_id,
                name=name,
                status='active',
                contact_email=f"contact@{name.lower().replace(' ', '')}.com",
                created_at=datetime.now().isoformat() + 'Z',
                updated_at=datetime.now().isoformat() + 'Z'
            )
            default_affiliates.append(affiliate)
        return default_affiliates

    def generate_click_events(self, num_clicks: int) -> List[ClickEvent]:
        """Generate realistic click events"""
        clicks = []
        
        # Generate clicks over the last 30 days
        end_date = datetime.now()
        start_date = end_date - timedelta(days=30)
        
        for i in range(num_clicks):
            # Random timestamp within the date range
            random_timestamp = start_date + timedelta(
                seconds=random.randint(0, int((end_date - start_date).total_seconds()))
            )
            
            # Select random campaign and affiliate
            campaign = random.choice(self.campaigns)
            affiliate = random.choice(self.affiliates)
            
            # Generate click ID
            click_id = f"click_{self.organization.organization_id}_{i+1:06d}"
            
            # Generate location data
            country = random.choice(self.countries)
            region = fake.state() if country == 'US' else fake.state()
            city = fake.city()
            
            # Generate referrer and landing page
            referrer_base = random.choice(self.referrer_sources)
            if 'search' in referrer_base:
                referrer_url = referrer_base + campaign.name.lower().replace(' ', '+')
            else:
                referrer_url = referrer_base + fake.word()
            
            landing_page_url = f"https://{self.organization.name.lower().replace(' ', '')}.com/{fake.word()}"
            
            # Generate sub parameters
            sources = ['google', 'facebook', 'instagram', 'youtube', 'tiktok', 'email', 'direct']
            mediums = ['cpc', 'social', 'organic', 'email', 'referral', 'video']
            campaigns_sub = ['summer', 'winter', 'holiday', 'sale', 'launch', 'promo']
            
            sub1 = f"source_{random.choice(sources)}"
            sub2 = f"medium_{random.choice(mediums)}"
            sub3 = f"campaign_{random.choice(campaigns_sub)}"
            
            # Determine if this click will convert (20% conversion rate)
            will_convert = random.random() < 0.2
            
            click = ClickEvent(
                organization_id=self.organization.organization_id,
                id=click_id,
                timestamp=random_timestamp.strftime('%Y-%m-%dT%H:%M:%SZ'),
                campaign_id=campaign.campaign_id,
                campaign_name=campaign.name,
                offer_id=f"offer_{campaign.campaign_id}",
                offer_name=f"{campaign.name} Offer",
                affiliate_id=affiliate.affiliate_id,
                affiliate_name=affiliate.name,
                ip_address=fake.ipv4(),
                user_agent=random.choice(self.user_agents),
                country=country,
                region=region,
                city=city,
                referrer_url=referrer_url,
                landing_page_url=landing_page_url,
                sub1=sub1,
                sub2=sub2,
                sub3=sub3,
                converted=will_convert,
                conversion_id=f"conv_{self.organization.organization_id}_{i+1:06d}" if will_convert else None
            )
            
            clicks.append(click)
        
        return clicks

    def generate_conversion_events(self, clicks: List[ClickEvent]) -> List[ConversionEvent]:
        """Generate conversion events based on clicks that converted"""
        conversions = []
        
        for click in clicks:
            if not click.converted or not click.conversion_id:
                continue
            
            # Conversion happens 1-60 minutes after click
            click_time = datetime.strptime(click.timestamp, '%Y-%m-%dT%H:%M:%SZ')
            conversion_time = click_time + timedelta(minutes=random.randint(1, 60))
            
            # Generate conversion value and payout
            conversion_value = round(random.uniform(50.0, 1000.0), 2)
            payout_percentage = random.uniform(0.05, 0.15)  # 5-15% payout
            payout = round(conversion_value * payout_percentage, 2)
            
            # Random status (80% approved, 15% pending, 5% rejected)
            status_rand = random.random()
            if status_rand < 0.8:
                status = 'approved'
            elif status_rand < 0.95:
                status = 'pending'
            else:
                status = 'rejected'
                payout = 0.0  # No payout for rejected conversions
            
            conversion = ConversionEvent(
                organization_id=self.organization.organization_id,
                id=click.conversion_id,
                timestamp=conversion_time.strftime('%Y-%m-%dT%H:%M:%SZ'),
                transaction_id=f"txn_{self.organization.name.lower().replace(' ', '_')}_{len(conversions)+1:06d}",
                campaign_id=click.campaign_id,
                campaign_name=click.campaign_name,
                offer_id=click.offer_id,
                offer_name=click.offer_name,
                status=status,
                payout=payout,
                currency=random.choice(self.currencies),
                affiliate_id=click.affiliate_id,
                affiliate_name=click.affiliate_name,
                click_id=click.id,
                conversion_value=conversion_value,
                sub1=click.sub1,
                sub2=click.sub2,
                sub3=click.sub3
            )
            
            conversions.append(conversion)
        
        return conversions

    def generate_performance_summary(self, clicks: List[ClickEvent], conversions: List[ConversionEvent]) -> List[PerformanceSummary]:
        """Generate performance summary data"""
        if not clicks:
            return []
        
        total_clicks = len(clicks)
        total_conversions = len(conversions)
        total_revenue = sum(c.conversion_value for c in conversions if c.status == 'approved')
        conversion_rate = (total_conversions / total_clicks * 100) if total_clicks > 0 else 0
        average_revenue = total_revenue / total_conversions if total_conversions > 0 else 0
        
        # Estimate impressions (clicks are typically 2-4% of impressions)
        click_through_rate = random.uniform(2.5, 3.5)
        total_impressions = int(total_clicks / (click_through_rate / 100))
        
        summary = PerformanceSummary(
            organization_id=self.organization.organization_id,
            total_clicks=total_clicks,
            total_conversions=total_conversions,
            total_revenue=round(total_revenue, 2),
            conversion_rate=round(conversion_rate, 2),
            average_revenue=round(average_revenue, 2),
            click_through_rate=round(click_through_rate, 1),
            total_impressions=total_impressions
        )
        
        return [summary]

    def generate_daily_performance_reports(self, clicks: List[ClickEvent], conversions: List[ConversionEvent]) -> List[DailyPerformanceReport]:
        """Generate daily performance reports"""
        if not clicks:
            return []
        
        # Group clicks by date and campaign
        daily_data = {}
        
        for click in clicks:
            date = click.timestamp[:10]  # Extract date part (YYYY-MM-DD)
            key = (date, click.campaign_id, click.campaign_name)
            
            if key not in daily_data:
                daily_data[key] = {
                    'clicks': 0,
                    'conversions': 0,
                    'revenue': 0.0,
                    'payouts': 0.0
                }
            
            daily_data[key]['clicks'] += 1
        
        # Add conversion data
        for conversion in conversions:
            if conversion.status != 'approved':
                continue
                
            date = conversion.timestamp[:10]
            key = (date, conversion.campaign_id, conversion.campaign_name)
            
            if key in daily_data:
                daily_data[key]['conversions'] += 1
                daily_data[key]['revenue'] += conversion.conversion_value
                daily_data[key]['payouts'] += conversion.payout
        
        # Generate daily reports
        reports = []
        for (date, campaign_id, campaign_name), data in daily_data.items():
            clicks_count = data['clicks']
            conversions_count = data['conversions']
            revenue = data['revenue']
            payouts = data['payouts']
            
            conversion_rate = (conversions_count / clicks_count * 100) if clicks_count > 0 else 0
            click_through_rate = random.uniform(2.8, 3.2)
            impressions = int(clicks_count / (click_through_rate / 100))
            
            report = DailyPerformanceReport(
                organization_id=self.organization.organization_id,
                date=date,
                campaign_id=campaign_id,
                campaign_name=campaign_name,
                clicks=clicks_count,
                impressions=impressions,
                conversions=conversions_count,
                revenue=round(revenue, 2),
                conversion_rate=round(conversion_rate, 1),
                click_through_rate=round(click_through_rate, 1),
                payouts=round(payouts, 2)
            )
            reports.append(report)
        
        return sorted(reports, key=lambda x: (x.date, x.campaign_id))

    def generate_campaign_performance(self, clicks: List[ClickEvent], conversions: List[ConversionEvent]) -> List[CampaignPerformance]:
        """Generate campaign performance data"""
        if not clicks:
            return []
        
        # Group data by campaign
        campaign_data = {}
        
        for click in clicks:
            if click.campaign_id not in campaign_data:
                campaign_data[click.campaign_id] = {
                    'name': click.campaign_name,
                    'clicks': 0,
                    'conversions': 0,
                    'revenue': 0.0
                }
            campaign_data[click.campaign_id]['clicks'] += 1
        
        # Add conversion data
        for conversion in conversions:
            if conversion.status != 'approved':
                continue
                
            if conversion.campaign_id in campaign_data:
                campaign_data[conversion.campaign_id]['conversions'] += 1
                campaign_data[conversion.campaign_id]['revenue'] += conversion.conversion_value
        
        # Generate campaign performance records
        performances = []
        statuses = ['active', 'paused', 'completed']
        
        for campaign_id, data in campaign_data.items():
            clicks_count = data['clicks']
            conversions_count = data['conversions']
            revenue = data['revenue']
            conversion_rate = (conversions_count / clicks_count * 100) if clicks_count > 0 else 0
            
            performance = CampaignPerformance(
                campaign_id=campaign_id,
                organization_id=self.organization.organization_id,
                name=data['name'],
                clicks=clicks_count,
                conversions=conversions_count,
                revenue=round(revenue, 2),
                conversion_rate=round(conversion_rate, 2),
                status=random.choice(statuses),
                priority=random.randint(1, 5),
                tier=random.randint(1, 5)
            )
            performances.append(performance)
        
        return performances

def save_to_csv(data: List[Any], filename: str, fieldnames: List[str]):
    """Save data to CSV file"""
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    
    with open(filename, 'w', newline='', encoding='utf-8') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        
        for item in data:
            if hasattr(item, '__dict__'):
                row = asdict(item)
            else:
                row = item
            
            # Handle None values
            for key, value in row.items():
                if value is None:
                    row[key] = ''
            
            writer.writerow(row)

def save_organization_info(org: Organization, campaigns: List[Campaign], 
                          advertisers: List[Advertiser], affiliates: List[Affiliate], 
                          output_dir: str):
    """Save organization information to JSON file"""
    org_info = {
        'organization': asdict(org),
        'campaigns': [asdict(c) for c in campaigns],
        'advertisers': [asdict(a) for a in advertisers],
        'affiliates': [asdict(a) for a in affiliates],
        'generated_at': datetime.now().isoformat()
    }
    
    info_file = os.path.join(output_dir, 'organization_info.json')
    with open(info_file, 'w', encoding='utf-8') as f:
        json.dump(org_info, f, indent=2, ensure_ascii=False)

def main():
    parser = argparse.ArgumentParser(description='Generate organization-specific mock data')
    parser.add_argument('--jwt-token', required=True, help='JWT token for authentication')
    parser.add_argument('--base-url', default='http://localhost:8080', help='Base URL of the API')
    parser.add_argument('--output-dir', default='./generated_mock_data', help='Output directory for generated data')
    parser.add_argument('--num-clicks', type=int, default=500, help='Number of click events to generate')
    parser.add_argument('--verbose', '-v', action='store_true', help='Enable verbose logging')
    
    args = parser.parse_args()
    
    if args.verbose:
        logger.setLevel(logging.DEBUG)
    
    try:
        # Initialize API client
        logger.info("Initializing API client...")
        api_client = APIClient(args.base_url, args.jwt_token)
        
        # Get user profile
        logger.info("Fetching user profile...")
        profile = api_client.get_user_profile()
        logger.info(f"User: {profile.email} (Role: {profile.role_name})")
        
        if not profile.organization_id:
            logger.error("User profile does not have an associated organization")
            return 1
        
        # Get organization details
        logger.info(f"Fetching organization details (ID: {profile.organization_id})...")
        organization = api_client.get_organization(profile.organization_id)
        logger.info(f"Organization: {organization.name} (Type: {organization.type})")
        
        # Get organization entities
        logger.info("Fetching organization campaigns...")
        campaigns = api_client.get_organization_campaigns(profile.organization_id)
        logger.info(f"Found {len(campaigns)} campaigns")
        
        logger.info("Fetching organization advertisers...")
        advertisers = api_client.get_organization_advertisers(profile.organization_id)
        logger.info(f"Found {len(advertisers)} advertisers")
        
        logger.info("Fetching organization affiliates...")
        affiliates = api_client.get_organization_affiliates(profile.organization_id)
        logger.info(f"Found {len(affiliates)} affiliates")
        
        # Create output directory
        org_output_dir = os.path.join(args.output_dir, f"org_{organization.organization_id}_{organization.name.replace(' ', '_')}")
        os.makedirs(org_output_dir, exist_ok=True)
        
        # Save organization information
        logger.info("Saving organization information...")
        save_organization_info(organization, campaigns, advertisers, affiliates, org_output_dir)
        
        # Generate mock data
        logger.info(f"Generating {args.num_clicks} click events...")
        generator = MockDataGenerator(organization, campaigns, advertisers, affiliates)
        
        clicks = generator.generate_click_events(args.num_clicks)
        logger.info(f"Generated {len(clicks)} click events")
        
        logger.info("Generating conversion events...")
        conversions = generator.generate_conversion_events(clicks)
        logger.info(f"Generated {len(conversions)} conversion events")
        
        # Save to CSV files
        logger.info("Saving click events to CSV...")
        clicks_file = os.path.join(org_output_dir, 'clicks_report.csv')
        click_fieldnames = [
            'organization_id', 'id', 'timestamp', 'campaign_id', 'campaign_name', 
            'offer_id', 'offer_name', 'affiliate_id', 'affiliate_name', 'ip_address', 
            'user_agent', 'country', 'region', 'city', 'referrer_url', 'landing_page_url', 
            'sub1', 'sub2', 'sub3', 'converted', 'conversion_id'
        ]
        save_to_csv(clicks, clicks_file, click_fieldnames)
        
        logger.info("Saving conversion events to CSV...")
        conversions_file = os.path.join(org_output_dir, 'conversions_report.csv')
        conversion_fieldnames = [
            'organization_id', 'id', 'timestamp', 'transaction_id', 'campaign_id', 
            'campaign_name', 'offer_id', 'offer_name', 'status', 'payout', 'currency', 
            'affiliate_id', 'affiliate_name', 'click_id', 'conversion_value', 
            'sub1', 'sub2', 'sub3'
        ]
        save_to_csv(conversions, conversions_file, conversion_fieldnames)
        
        # Generate dashboard and reporting data
        logger.info("Generating dashboard and reporting data...")
        
        # Performance summary
        logger.info("Generating performance summary...")
        performance_summaries = generator.generate_performance_summary(clicks, conversions)
        performance_summary_file = os.path.join(org_output_dir, 'performance_summary.csv')
        performance_summary_fieldnames = [
            'organization_id', 'total_clicks', 'total_conversions', 'total_revenue',
            'conversion_rate', 'average_revenue', 'click_through_rate', 'total_impressions'
        ]
        save_to_csv(performance_summaries, performance_summary_file, performance_summary_fieldnames)
        
        # Daily performance reports
        logger.info("Generating daily performance reports...")
        daily_reports = generator.generate_daily_performance_reports(clicks, conversions)
        daily_reports_file = os.path.join(org_output_dir, 'daily_performance_report.csv')
        daily_reports_fieldnames = [
            'organization_id', 'date', 'campaign_id', 'campaign_name', 'clicks',
            'impressions', 'conversions', 'revenue', 'conversion_rate', 
            'click_through_rate', 'payouts'
        ]
        save_to_csv(daily_reports, daily_reports_file, daily_reports_fieldnames)
        
        # Campaign performance
        logger.info("Generating campaign performance data...")
        campaign_performances = generator.generate_campaign_performance(clicks, conversions)
        campaign_performance_file = os.path.join(org_output_dir, 'campaign_performance.csv')
        campaign_performance_fieldnames = [
            'campaign_id', 'organization_id', 'name', 'clicks', 'conversions',
            'revenue', 'conversion_rate', 'status', 'priority', 'tier'
        ]
        save_to_csv(campaign_performances, campaign_performance_file, campaign_performance_fieldnames)
        
        # Generate summary
        total_revenue = sum(c.conversion_value for c in conversions if c.status == 'approved')
        total_payout = sum(c.payout for c in conversions if c.status == 'approved')
        conversion_rate = (len(conversions) / len(clicks)) * 100 if clicks else 0
        
        logger.info("=" * 60)
        logger.info("GENERATION SUMMARY")
        logger.info("=" * 60)
        logger.info(f"Organization: {organization.name} (ID: {organization.organization_id})")
        logger.info(f"Output Directory: {org_output_dir}")
        logger.info(f"Total Clicks: {len(clicks)}")
        logger.info(f"Total Conversions: {len(conversions)}")
        logger.info(f"Conversion Rate: {conversion_rate:.2f}%")
        logger.info(f"Total Revenue: ${total_revenue:.2f}")
        logger.info(f"Total Payout: ${total_payout:.2f}")
        logger.info(f"Campaigns Used: {len(campaigns)}")
        logger.info(f"Affiliates Used: {len(affiliates)}")
        logger.info("")
        logger.info("Generated Files:")
        logger.info("- clicks_report.csv (raw click events)")
        logger.info("- conversions_report.csv (raw conversion events)")
        logger.info("- performance_summary.csv (overall performance metrics)")
        logger.info("- daily_performance_report.csv (daily campaign performance)")
        logger.info("- campaign_performance.csv (campaign-level metrics)")
        logger.info("- organization_info.json (metadata and configuration)")
        logger.info("=" * 60)
        
        return 0
        
    except requests.exceptions.HTTPError as e:
        logger.error(f"HTTP Error: {e}")
        logger.error(f"Response: {e.response.text if e.response else 'No response'}")
        return 1
    except Exception as e:
        logger.error(f"Error: {e}")
        return 1

if __name__ == '__main__':
    exit(main())