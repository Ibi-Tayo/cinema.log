-- +goose Up
-- +goose StatementBegin
CREATE TABLE film_recommendation (
    film_recommendation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    external_film_id UUID NOT NULL,
    has_seen BOOLEAN NOT NULL,
    has_been_recommended BOOLEAN NOT NULL,
    recommendations_generated BOOLEAN NOT NULL,
);
ALTER TABLE film_recommendation 
ADD CONSTRAINT pk_film_recommendation PRIMARY KEY (film_recommendation_id);

ALTER TABLE film_recommendation 
ADD CONSTRAINT fk_film_recommendation_users_user_id 
FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;

ALTER TABLE film_recommendation 
ADD CONSTRAINT fk_film_recommendation_films_external_film_id 
FOREIGN KEY (external_film_id) REFERENCES films (external_id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS film_recommendation CASCADE;
-- +goose StatementEnd
