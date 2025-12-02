---
name: velocity
description: Generate weekly velocity report from Linear
allowed_tools: []
---

# Weekly Velocity Report

Use Linear MCP to analyze the past week's development velocity and generate a comprehensive report.

## Report Sections to Generate

### 1. Throughput Analysis
Query Linear for:
- Issues completed this week vs. last week
- Issues by type (bug, feature, docs, chore)
- Issues by team (Core, Frontend, Infra)

Format as:
```
| Metric | This Week | Last Week | Change |
|--------|-----------|-----------|--------|
| Completed | X | Y | +/-Z% |
| Bugs Fixed | X | Y | +/-Z% |
| Features Shipped | X | Y | +/-Z% |
```

### 2. Cycle Time
Calculate average time from:
- Created → In Progress (pickup time)
- In Progress → Done (completion time)
- Total cycle time

Identify outliers (issues that took unusually long).

### 3. Work Distribution
- Issues by assignee
- Issues by priority level
- Blocked time analysis

### 4. Sprint Health
- Burndown progress (if using cycles)
- Scope changes (issues added/removed mid-sprint)
- Carryover from previous sprint

### 5. Recommendations
Based on the data, suggest:
- Process improvements
- Bottleneck areas
- Resource allocation adjustments
- Upcoming risks

## Output Format

Generate a markdown report suitable for:
1. Pasting into Slack/Discord
2. Adding to CLAUDE.md
3. Sharing in team meetings

Include:
- Executive summary (3 bullet points)
- Detailed metrics table
- Trend analysis
- Action items

## Example Output

```markdown
## Velocity Report: Week of Dec 2, 2025

### TL;DR
- Shipped 12 issues (up 20% from last week)
- Average cycle time: 2.3 days (improved from 3.1)
- One P1 bug still in progress >4 days - needs attention

### Throughput
| Type | Completed | WoW Change |
|------|-----------|------------|
| Features | 5 | +25% |
| Bugs | 4 | -10% |
| Chores | 3 | +50% |

### Cycle Time
- Avg pickup: 4 hours
- Avg completion: 2.1 days
- Outlier: DRI-234 (6 days) - blocked on external API

### Action Items
1. Review DRI-234 blocking issue
2. Consider pairing on Frontend backlog (growing)
3. Schedule tech debt sprint for January
```

## Notes

- If Linear data is limited, note what's missing
- Compare to baseline if available
- Flag any data quality issues (unlabeled issues, etc.)
