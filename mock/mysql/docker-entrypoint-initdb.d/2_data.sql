USE example;

START TRANSACTION;

INSERT INTO user (id, name) VALUES
                                (1, "Hoge"),
                                (2, "Fuga");


INSERT INTO message (id, user_id, message) VALUES
                                               (1,1,"test message id 1"),
                                               (2,1,"test message id 2"),
                                               (3,2,"test message id 3"),
                                               (4,2,"test message id 4");

COMMIT;
