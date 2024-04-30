--THDS:Up
--Implement up part of migration here 
CREATE TABLE dbo.TestTable
(
  ID int IDENTITY(1, 1) PRIMARY KEY
);

CREATE TABLE dbo.TestTable2
(
  ID int IDENTITY(1, 1) PRIMARY KEY
);
GO

--THDS:Down
-- Drop the table 'TestTable' in schema 'dbo'
IF EXISTS (
  SELECT *
FROM sys.tables
  JOIN sys.schemas
  ON sys.tables.schema_id = sys.schemas.schema_id
WHERE sys.schemas.name = N'dbo'
  AND sys.tables.name = N'TestTable'
)
  DROP TABLE dbo.TestTable
GO

-- Drop the table 'TestTable2' in schema 'dbo'
IF EXISTS (
  SELECT *
FROM sys.tables
  JOIN sys.schemas
  ON sys.tables.schema_id = sys.schemas.schema_id
WHERE sys.schemas.name = N'dbo'
  AND sys.tables.name = N'TestTable2'
)
  DROP TABLE dbo.TestTable2
GO
