CREATE DATABASE tasks_db;

CREATE TABLE IF NOT EXISTS tasks
(
    id               BIGSERIAL PRIMARY KEY,
    user_id          integer REFERENCES users (id) ON DELETE CASCADE,
    task_name        VARCHAR(30)  NOT NULL,
    description      VARCHAR(255) NOT NULL,
    is_done          BOOLEAN      NOT NULL,
    added            TIMESTAMP    NOT NULL DEFAULT now(),
    deadline         TIMESTAMP,
    done_at          TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users
(
    id                  SERIAL PRIMARY KEY,
    full_name           VARCHAR(30) NOT NULL,
    user_name           VARCHAR(30) NOT NULL,
    password_hash       VARCHAR(255) NOT NULL
--     date_of_Birth DATE NOT NULL
);


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

