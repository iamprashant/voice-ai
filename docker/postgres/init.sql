-- Create databases for each service
CREATE DATABASE IF NOT EXISTS web_db;
CREATE DATABASE IF NOT EXISTS assistant_db;
CREATE DATABASE IF NOT EXISTS integration_db;
CREATE DATABASE IF NOT EXISTS endpoint_db;

-- Grant privileges to rapida_user on all databases
GRANT ALL PRIVILEGES ON DATABASE web_db TO rapida_user;
GRANT ALL PRIVILEGES ON DATABASE assistant_db TO rapida_user;
GRANT ALL PRIVILEGES ON DATABASE integration_db TO rapida_user;
GRANT ALL PRIVILEGES ON DATABASE endpoint_db TO rapida_user;
