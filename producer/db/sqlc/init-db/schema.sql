CREATE TYPE task_state AS ENUM ('pending', 'in_progress', 'completed', 'failed');

CREATE TABLE tasks (
    ID INT PRIMARY KEY,
    Type INT CHECK (Type BETWEEN 0 AND 9) NOT NULL,
    Value INT CHECK (Value BETWEEN 0 AND 99) NOT NULL,
    State task_state NOT NULL,
    CreationTime FLOAT NOT NULL,
    LastUpdateTime FLOAT NOT NULL
);