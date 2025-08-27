#!/usr/bin/env python3
"""
Helper script to add sample campaigns, advertisers, and affiliates to an organization.
"""

import psycopg2
import argparse
import logging
from datetime import datetime

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def connect_to_db():
    """Connect to PostgreSQL database"""
    try:
        conn = psycopg2.connect(
            host="localhost",
            database="affiliate_platform",
            user="postgres",
            password="postgres",
            port="5432"
        )
        return conn
    except Exception as e:
        logger.error(f"Failed to connect to database: {e}")
        return None

def add_sample_data(org_id):
    """Add sample campaigns, advertisers, and affiliates to organization"""
    conn = connect_to_db()
    if not conn:
        return False
    
    try:
        cursor = conn.cursor()
        
        # Add sample advertiser
        cursor.execute("""
            INSERT INTO advertisers (organization_id, name, status, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s)
            RETURNING advertiser_id
        """, (org_id, "Dyson Advertiser", "active", datetime.now(), datetime.now()))
        
        advertiser_id = cursor.fetchone()[0]
        logger.info(f"Created advertiser: Dyson Advertiser (ID: {advertiser_id})")
        
        # Add sample affiliate
        cursor.execute("""
            INSERT INTO affiliates (organization_id, name, status, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s)
            RETURNING affiliate_id
        """, (org_id, "Dyson Affiliate", "active", datetime.now(), datetime.now()))
        
        affiliate_id = cursor.fetchone()[0]
        logger.info(f"Created affiliate: Dyson Affiliate (ID: {affiliate_id})")
        
        # Add sample campaign
        cursor.execute("""
            INSERT INTO campaigns (organization_id, advertiser_id, name, status, visibility, destination_url, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
            RETURNING campaign_id
        """, (org_id, advertiser_id, "Dyson V15 Campaign", "active", "public", "https://dyson.com/v15", datetime.now(), datetime.now()))
        
        campaign_id = cursor.fetchone()[0]
        logger.info(f"Created campaign: Dyson V15 Campaign (ID: {campaign_id})")
        
        conn.commit()
        cursor.close()
        conn.close()
        
        logger.info(f"Successfully added sample data to organization {org_id}")
        return True
        
    except Exception as e:
        logger.error(f"Failed to add sample data: {e}")
        conn.rollback()
        conn.close()
        return False

def main():
    parser = argparse.ArgumentParser(description='Add sample data to organization')
    parser.add_argument('--org-id', type=int, required=True, help='Organization ID')
    parser.add_argument('--verbose', '-v', action='store_true', help='Enable verbose logging')
    
    args = parser.parse_args()
    
    if args.verbose:
        logger.setLevel(logging.DEBUG)
    
    success = add_sample_data(args.org_id)
    
    if success:
        logger.info("Sample data successfully added")
        return 0
    else:
        logger.error("Failed to add sample data")
        return 1

if __name__ == "__main__":
    exit(main())