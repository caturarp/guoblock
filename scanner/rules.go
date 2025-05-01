package scanner

import "regexp"

type Rule struct {
    Name     string
    Pattern  *regexp.Regexp
    Priority int // Higher priority rules are checked first
}

var Rules = []Rule{
    {
        Name:     "AWS Access Key",
        Pattern:  regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
        Priority: 100,
    },
    {
        Name:     "AWS Secret Key",
        Pattern:  regexp.MustCompile(`(?i)aws(.{0,20})?(secret|key)['"]?\s*[:=]\s*['"][0-9a-zA-Z\/+]{40}['"]`),
        Priority: 100,
    },
    {
        Name:     "Slack Token",
        Pattern:  regexp.MustCompile(`(?i)xox[bap]-[0-9a-zA-Z]{10,48}`),
        Priority: 90,
    },
    {
        Name:     "Discord Webhook",
        Pattern:  regexp.MustCompile(`https:\/\/discord\.com\/api\/webhooks\/[0-9]{18,19}\/[a-zA-Z0-9_-]+`),
        Priority: 90,
    },
    {
        Name:     "GCP Service Account JSON",
        Pattern:  regexp.MustCompile(`"type": "service_account".*"project_id": "[a-z0-9-]+"`),
        Priority: 90,
    },
    {
        Name:     "Generic API Key",
        Pattern:  regexp.MustCompile(`(?i)(api|token|key)['"]?\s*[:=]\s*['"][0-9a-zA-Z_\-]{16,}['"]`),
        Priority: 10, // Lowest priority as it's the most generic
    },
}
