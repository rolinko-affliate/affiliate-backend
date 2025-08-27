# Mock Data Generation Script - Usage Guide

## Overview
This script generates realistic mock data for affiliate marketing organizations using JWT tokens for authentication. It queries the API to fetch organization details and generates hundreds of click and conversion events with proper business logic.

## Prerequisites
1. **API Server Running**: Ensure the affiliate platform API server is running on port 18080
2. **Valid JWT Token**: You need a valid JWT token for authentication
3. **User in Database**: The user associated with the JWT token must exist in the database

## Quick Start

### 1. Generate Mock Data with Your JWT Token
```bash
python generate_org_mock_data.py --jwt-token "YOUR_JWT_TOKEN_HERE" --num-clicks 500
```

### 2. Add User to Database (if needed)
If your JWT token user doesn't exist in the database:
```bash
python add_user_to_db.py --jwt-token "YOUR_JWT_TOKEN_HERE" --org-name "Your Organization Name"
```

### 3. Add Sample Data to Organization (if needed)
If your organization has no campaigns/advertisers/affiliates:
```bash
python add_sample_data.py --org-id YOUR_ORG_ID
```

## Script Options

### generate_org_mock_data.py
- `--jwt-token`: Your JWT authentication token (required)
- `--num-clicks`: Number of click events to generate (default: 500)
- `--base-url`: API base URL (default: http://localhost:18080)
- `--output-dir`: Output directory for generated files (default: ./generated_mock_data)

### add_user_to_db.py
- `--jwt-token`: JWT token containing user information (required)
- `--org-name`: Organization name to create/associate with user
- `--verbose`: Enable detailed logging

### add_sample_data.py
- `--org-id`: Organization ID to add sample data to (required)
- `--verbose`: Enable detailed logging

## Generated Files

The script creates a folder structure like `generated_mock_data/org_2_Dyson/` containing:

1. **clicks_report.csv** - Raw click event data with tracking parameters
2. **conversions_report.csv** - Conversion events with revenue and payout data
3. **performance_summary.csv** - Overall performance metrics
4. **daily_performance_report.csv** - Daily campaign performance breakdown
5. **campaign_performance.csv** - Campaign-level performance metrics
6. **organization_info.json** - Organization metadata and configuration

## Sample Output
```
Organization: Dyson (ID: 2)
Total Clicks: 1000
Total Conversions: 213
Conversion Rate: 21.30%
Total Revenue: $87,708.39
Total Payout: $9,425.39
```

## Features

### Realistic Data Generation
- **Geographic Distribution**: Clicks from various countries and regions
- **Device Variety**: Mobile, desktop, and tablet traffic
- **Time Distribution**: Events spread over 30-day period
- **Conversion Logic**: Realistic conversion rates (15-25%)
- **Revenue Modeling**: Variable conversion values and payouts
- **Attribution**: Proper click-to-conversion attribution

### Business Logic
- **Campaign Performance**: Tracks performance by campaign
- **Affiliate Attribution**: Associates clicks/conversions with affiliates
- **Currency Support**: Multi-currency conversion values
- **Status Tracking**: Approved, pending, and rejected conversions
- **UTM Parameters**: Realistic tracking parameters (source, medium, campaign)

### Data Quality
- **Consistent IDs**: Proper ID generation and referencing
- **Timestamp Accuracy**: Realistic time distributions
- **Data Validation**: Ensures data integrity and relationships
- **CSV Formatting**: Clean, importable CSV files

## Troubleshooting

### Common Issues

1. **"User not found" Error**
   - Solution: Run `add_user_to_db.py` to add your user to the database

2. **"No campaigns/advertisers/affiliates found"**
   - Solution: Run `add_sample_data.py` to add sample data to your organization

3. **"Connection refused" Error**
   - Solution: Ensure the API server is running on port 18080

4. **"Invalid JWT token" Error**
   - Solution: Check that your JWT token is valid and not expired

### Database Connection
The scripts connect to PostgreSQL with these default settings:
- Host: localhost
- Database: affiliate_platform
- User: postgres
- Password: postgres
- Port: 5432

## Example Workflow

1. **Start API Server** (if not running)
2. **Add User to Database**:
   ```bash
   python add_user_to_db.py --jwt-token "YOUR_TOKEN" --org-name "Dyson"
   ```
3. **Add Sample Data** (if organization is empty):
   ```bash
   python add_sample_data.py --org-id 2
   ```
4. **Generate Mock Data**:
   ```bash
   python generate_org_mock_data.py --jwt-token "YOUR_TOKEN" --num-clicks 1000
   ```

## Data Schema

### Click Events
- Organization and campaign attribution
- Affiliate tracking
- Geographic and device information
- UTM parameters and referrer data
- Conversion status and linking

### Conversion Events
- Revenue and payout calculations
- Multi-currency support
- Conversion status (approved/pending/rejected)
- Attribution to original clicks
- Transaction IDs and metadata

This script provides a comprehensive solution for generating realistic affiliate marketing data for testing, development, and demonstration purposes.