# Service Documentation Package

This directory contains the comprehensive service documentation package for the Affiliate Backend Platform, organized according to ITIL4/DevOps/SRE best practices.

## Document Structure

### Strategic Layer
Documents for C-level executives and product owners (non-technical language):

1. **[Service Catalogue Entry](01-service-catalogue-entry.md)** - One-page summary for business users
2. **[Service-Level Agreement (SLA)](02-service-level-agreement.md)** - Formal commitment & metrics
3. **[Operating Model Overview](03-operating-model-overview.md)** - "Who does what" at a glance

### Service Layer
Documents for architects and auditors (precise diagrams, references):

4. **[System Architecture Overview](04-system-architecture-overview.md)** - 30-minute orientation for newcomers
5. **[Configuration Baseline (CMDB / IaC)](05-configuration-baseline.md)** - Source-of-truth for infrastructure
6. **[Security & Compliance Guide](06-security-compliance-guide.md)** - Show regulators we're covered
7. **[BCP / DR Plan](07-bcp-dr-plan.md)** - "If it breaks, how fast can we be back?"

### Operational Layer
Documents for on-call engineers (command-ready, copy-paste):

8. **[Runbook](08-runbook.md)** - Step-by-step actions for on-call
9. **[Monitoring & Alerting Playbook](09-monitoring-alerting-playbook.md)** - Keep eyes on SLOs

## Document Standards

- **Version Control**: Each document includes owner and next review date (≤ 6 months)
- **Semantic Versioning**: vMAJOR.MINOR format
- **Single Source of Truth**: Markdown in Git, diagrams in PlantUML/draw.io
- **Measurable Content**: All SLOs with unit + target
- **Tested Procedures**: Runbook steps validated in staging

## Quick Navigation

| Layer | Document | Owner | Last Updated | Next Review |
|-------|----------|-------|--------------|-------------|
| Strategic | [Service Catalogue Entry](01-service-catalogue-entry.md) | Product Owner | 2025-08-05 | 2026-02-05 |
| Strategic | [Service-Level Agreement](02-service-level-agreement.md) | Service Manager | 2025-08-05 | 2026-02-05 |
| Strategic | [Operating Model Overview](03-operating-model-overview.md) | Operations Manager | 2025-08-05 | 2026-02-05 |
| Service | [System Architecture Overview](04-system-architecture-overview.md) | Lead Architect | 2025-08-05 | 2026-02-05 |
| Service | [Configuration Baseline](05-configuration-baseline.md) | DevOps Engineer | 2025-08-05 | 2026-02-05 |
| Service | [Security & Compliance Guide](06-security-compliance-guide.md) | Security Officer | 2025-08-05 | 2026-02-05 |
| Service | [BCP / DR Plan](07-bcp-dr-plan.md) | SRE Lead | 2025-08-05 | 2026-02-05 |
| Operational | [Runbook](08-runbook.md) | On-Call Team Lead | 2025-08-05 | 2026-02-05 |
| Operational | [Monitoring & Alerting Playbook](09-monitoring-alerting-playbook.md) | SRE Team | 2025-08-05 | 2026-02-05 |

---

**Package Version**: v1.0  
**Last Updated**: 2025-08-05  
**Next Package Review**: 2026-02-05