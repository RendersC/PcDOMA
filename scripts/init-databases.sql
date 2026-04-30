-- Create separate databases for each service (database-per-service pattern)
CREATE DATABASE pcrentalauth;
CREATE DATABASE pcrentalbooking;
CREATE DATABASE pcrentalpayment;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE pcrentalauth TO pcdoma;
GRANT ALL PRIVILEGES ON DATABASE pcrentalbooking TO pcdoma;
GRANT ALL PRIVILEGES ON DATABASE pcrentalpayment TO pcdoma;
