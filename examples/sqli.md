# Basic SQL Syntax
(select|union|insert|update|delete|replace|alter|create|drop|truncate)\s+.*?(from|into|table|database)
(select|union|insert|update|delete|replace|alter|create|drop|truncate)\s*\(.*?\)
(select|union|insert|update|delete|replace|alter|create|drop|truncate)\s+.*?\s*\(.*?\)
'or\s+1=1\s*--\s*
'or\s+'[^']+'='\1\s*--\s*
'[^']*'\s*=\s*'[^']*'\s*--\s*
[^a-zA-Z0-9]\s*(or|and)\s+(\d+)=\2
[^a-zA-Z0-9]\s*(or|and)\s+'[^']*'='[^']*'

# SQL Comments
--\s+
#\s+
;\s*--
;\s*#
;\s*/\*\s*\w*\s*\*/\s*

# SQL Functions and Operators
(select|union|sleep|benchmark|extractvalue|updatexml|load_file|floor|rand|md5|sha1|if)\s*\(
(select|union|sleep|benchmark|extractvalue|updatexml|load_file|floor|rand|md5|sha1|if)\s*[\(\w]+
case\s+when\s+.*?\s+then\s+.*?\s+end
if\((.*?)\)\s*{.*?}\s*else\s*{.*?}

# SQL Boolean Expressions
'(\d+)=\1
'[^']*'='[^']*'
'(\d+)\s+(or|and)\s+\d
'(\d+)\s+(or|and)\s+\d\s*--\s*
