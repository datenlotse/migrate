--THDS:Up
--Implement up part of migration here 
-- Create a new table called 'test3' in schema 'dbo'
-- Drop the table if it already exists
IF OBJECT_ID('dbo.test3', 'U') IS NOT NULL
DROP TABLE dbo.test3
GO
-- Create the table in the specified schema
CREATE TABLE dbo.test3
(
  test3Id INT NOT NULL PRIMARY KEY,
  -- primary key column
);
GO

--THDS:Down
--Implement down part of migration here
-- Drop the table 'test3' in schema 'dbo'
IF EXISTS (
  SELECT *
FROM sys.tables
  JOIN sys.schemas
  ON sys.tables.schema_id = sys.schemas.schema_id
WHERE sys.schemas.name = N'dbo'
  AND sys.tables.name = N'test3'
)
  DROP TABLE dbo.test3
GO