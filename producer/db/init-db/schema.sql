CREATE TYPE task_state AS ENUM ('pending', 'in_progress', 'completed', 'failed');

CREATE TABLE tasks (
    ID SERIAL PRIMARY KEY,
    Type INT CHECK (Type BETWEEN 0 AND 9),
    Value INT CHECK (Value BETWEEN 0 AND 99),
    State task_state,
    CreationTime FLOAT,
    LastUpdateTime FLOAT
);