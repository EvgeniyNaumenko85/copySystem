CREATE DATABASE files_db;

CREATE TABLE IF NOT EXISTS users
(
    id                  BIGSERIAL PRIMARY KEY,
    user_name           VARCHAR(30) NOT NULL UNIQUE,
    email               VARCHAR(30) NOT NULL UNIQUE,
    password_hash       VARCHAR(255) NOT NULL,
    user_role           VARCHAR(30) NOT NULL,
    file_size_lim       INT DEFAULT 20,
);

CREATE TABLE IF NOT EXISTS files
(
    id               BIGSERIAL PRIMARY KEY,
    user_id          INTEGER REFERENCES users (id) ON DELETE CASCADE,
    file_name        VARCHAR(30)  NOT NULL,
    extension        VARCHAR(10),
    file_path        VARCHAR(255) NOT NULL UNIQUE,
    file_size        INTEGER,
    deleted          BOOLEAN      NOT NULL,
    added            TIMESTAMP    NOT NULL DEFAULT now(),
);

CREATE TABLE IF NOT EXISTS accesses
(
    id               BIGSERIAL PRIMARY KEY,
    user_id          INTEGER REFERENCES users (id) ON DELETE CASCADE,
    file_id          INTEGER REFERENCES files (id) ON DELETE CASCADE,
);


/*
-- Creating  function getUndoneTasksByUserId:
CREATE OR REPLACE FUNCTION getUndoneTasksByUserId(user_id_arg integer)
    RETURNS TABLE (id bigint, user_id int, task_name varchar(30), description varchar(255), is_done boolean, added timestamp, deadline timestamp)
    LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN QUERY
        SELECT tasks.id, tasks.user_id, tasks.task_name, tasks.description, tasks.is_done, tasks.added, tasks.deadline
        FROM tasks
        WHERE tasks.user_id = user_id_arg AND tasks.is_done = false;
END;
$$;

SELECT EXISTS (
    SELECT 1
    FROM pg_proc
    WHERE proname = 'getundonetasksbyuserid'
);

SELECT current_database();

SELECT * FROM getUndoneTasksByUserId(2);


-- Creating  function reassign_task:
CREATE OR REPLACE PROCEDURE reassign_task(
    old_user_id INTEGER,
    new_user_id INTEGER,
    task_id BIGINT
)
    LANGUAGE plpgsql
AS $$
BEGIN
    -- Обновляем запись задачи, переназначая ее на нового пользователя
    UPDATE tasks
    SET user_id = new_user_id
    WHERE id = task_id AND user_id = old_user_id;

    -- Проверяем, была ли произведена переназначение задачи
    IF FOUND THEN
        RAISE NOTICE 'Задача % переназначена с пользователя % на пользователя %',
            task_id, old_user_id, new_user_id;
    ELSE
        RAISE EXCEPTION 'Задача % не найдена для пользователя %', task_id, old_user_id;
    END IF;
END;
$$;

SELECT EXISTS (
    SELECT 1
    FROM pg_proc
    WHERE proname = 'reassign_task'
);

CALL reassign_task(1, 21, 1);


-- Creating  function get_overdue_tasks:
CREATE OR REPLACE FUNCTION get_overdue_tasks(user_indent INT)
    RETURNS TABLE (
                      id BIGINT,
                      user_id INT,
                      task_name VARCHAR(30),
                      is_done BOOLEAN,
                      description VARCHAR(255),
                      added TIMESTAMP,
                      deadline TIMESTAMP
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.user_id, t.task_name, t.is_done, t.description, t.added, t.deadline
        FROM tasks t
        WHERE t.user_id = user_indent AND t.deadline < NOW() AND t.is_done = FALSE;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION get_overdue_tasks(integer);

SELECT EXISTS (
    SELECT 1
    FROM pg_proc
    WHERE proname = 'get_overdue_tasks'
);

SELECT * FROM get_overdue_tasks(2);

-- ========================>>
*/
