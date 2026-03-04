package llm

import "fmt"

const initialOutreachTemplate = `You are a business development executive for Taraka International, a merchant exporter based in Mumbai, India. 
We export FMCG Foods, Superfoods, Nutraceuticals, Horticulture products (Cocopeat), and Eco-Packaging to global markets.

Write a professional, concise cold email to %s at %s (%s) in %s.
The email should:
1. Introduce Taraka International briefly
2. Mention our products relevant to their business
3. Propose a business partnership
4. Have a clear call to action
5. Be under 200 words
6. Sound natural and human, not like AI

Company: %s
Website: %s

Return ONLY the email body, no subject line.`

const followUp1Template = `Write a brief, friendly follow-up email to %s at %s.
This is the first follow-up to our initial outreach about Taraka International's export products.
Keep it under 100 words. Reference our previous email. Sound human and natural.
Return ONLY the email body.`

const followUp2Template = `Write a second follow-up email to %s at %s.
This is our second attempt to connect. Be brief, value-focused.
Mention one specific product that could benefit their business.
Under 80 words. Sound natural.
Return ONLY the email body.`

const followUp3Template = `Write a final follow-up email to %s at %s.
This is our last attempt. Keep it very short (under 60 words).
Be respectful of their time. Leave the door open for future contact.
Return ONLY the email body.`

const replyResponseTemplate = `You are a business development executive for Taraka International.
A potential client has replied to our outreach email:

---
%s
---

Write a professional, helpful response that:
1. Acknowledges their reply
2. Answers any questions they have
3. Moves the conversation forward
4. Is under 150 words
5. Sounds natural and human

Return ONLY the email body.`

const businessMatchTemplate = `Analyze if this business would be a good potential buyer for Taraka International's exports.

Company: %s
Website: %s
Country: %s
Business Type: %s

Taraka International exports: FMCG Foods, Superfoods, Nutraceuticals, Horticulture (Cocopeat), Eco-Packaging
Target markets: UAE, UK, USA, Australia, Singapore, Europe, Canada, New Zealand, Japan, South Korea, Africa

Rate the match from 0-100 and explain briefly.
Return JSON: {"score": 75, "reason": "explanation"}`

const nlToSQLTemplate = `Convert this natural language query to SQL for a leads database.

Schema:
- leads (id, name, email, phone, company, website, country, source, business_type, score, status, emails_sent, follow_ups_sent, reply_received, notes, created_at)
- status values: 'raw', 'verified', 'emailed', 'follow_up_1', 'follow_up_2', 'follow_up_3', 'replied', 'converted', 'rejected'

Query: %s

Return ONLY the SQL query, nothing else.`

func InitialOutreachPrompt(name, company, businessType, country, website, notes string) string {
	return fmt.Sprintf(initialOutreachTemplate, name, company, businessType, country, company, website)
}

func FollowUp1Prompt(name, company string) string {
	return fmt.Sprintf(followUp1Template, name, company)
}

func FollowUp2Prompt(name, company string) string {
	return fmt.Sprintf(followUp2Template, name, company)
}

func FollowUp3Prompt(name, company string) string {
	return fmt.Sprintf(followUp3Template, name, company)
}

func ReplyResponsePrompt(replyText string) string {
	return fmt.Sprintf(replyResponseTemplate, replyText)
}

func BusinessMatchPrompt(company, website, country, businessType string) string {
	return fmt.Sprintf(businessMatchTemplate, company, website, country, businessType)
}

func NLToSQLPrompt(query string) string {
	return fmt.Sprintf(nlToSQLTemplate, query)
}
