package memory

const createLeadsTable = `
CREATE TABLE IF NOT EXISTS leads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT UNIQUE,
    phone TEXT,
    company TEXT,
    website TEXT,
    country TEXT,
    source TEXT,
    business_type TEXT,
    score INTEGER DEFAULT 0,
    status TEXT DEFAULT 'raw',
    gmail_thread_id TEXT,
    gmail_message_id TEXT,
    emails_sent INTEGER DEFAULT 0,
    follow_ups_sent INTEGER DEFAULT 0,
    last_emailed_at DATETIME,
    next_follow_up DATETIME,
    reply_received BOOLEAN DEFAULT 0,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createDailyLogTable = `
CREATE TABLE IF NOT EXISTS daily_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE UNIQUE,
    leads_found INTEGER DEFAULT 0,
    leads_verified INTEGER DEFAULT 0,
    emails_sent INTEGER DEFAULT 0,
    follow_ups_sent INTEGER DEFAULT 0,
    replies_received INTEGER DEFAULT 0,
    errors TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createCompanyKnowledgeTable = `
CREATE TABLE IF NOT EXISTS company_knowledge (
    key TEXT PRIMARY KEY,
    value TEXT
);`

const createCrawlHistoryTable = `
CREATE TABLE IF NOT EXISTS crawl_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT,
    source TEXT,
    country TEXT,
    results_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createConversationsTable = `
CREATE TABLE IF NOT EXISTS conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_query TEXT,
    bot_response TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
