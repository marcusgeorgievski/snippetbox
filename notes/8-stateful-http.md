# Stateful HTTP

- Session managers
- Use sesions to safely and securely share data between requests for a user
- Customize session behaviour (timeouts, cookies)

## Session Manager + Working with session data

1. Sessions table - `token`, `data`, `expiry`
2. Add session manager to `application` struct
3. Register session manager loader/saver to middleware, automate key retrieval in `newTemplateData`
4. Put, Get, Pop operations on keys

**How it works**

- Session manager adds session token to cookies, creates record in db with same token + blob data
