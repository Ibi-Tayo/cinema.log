-- +goose Up
-- +goose StatementBegin
CREATE TABLE film_graph_nodes (
    user_id UUID NOT NULL,
    external_film_id INT NOT NULL,
    title VARCHAR(255) NOT NULL
);
ALTER TABLE film_graph_nodes 
ADD CONSTRAINT pk_film_graph_nodes PRIMARY KEY (user_id, external_film_id);
ALTER TABLE film_graph_nodes 
ADD CONSTRAINT fk_film_graph_nodes_users_user_id 
FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;
ALTER TABLE film_graph_nodes 
ADD CONSTRAINT fk_film_graph_nodes_films_external_film_id 
FOREIGN KEY (external_film_id) REFERENCES films (external_id) ON DELETE CASCADE; 


CREATE TABLE film_graph_edges (
    user_id UUID NOT NULL,
    edge_id UUID NOT NULL,
    from_film_id INT NOT NULL,
    to_film_id INT NOT NULL
);
ALTER TABLE film_graph_edges 
ADD CONSTRAINT pk_film_graph_edges PRIMARY KEY (edge_id);
ALTER TABLE film_graph_edges 
ADD CONSTRAINT fk_film_graph_edges_users_user_id 
FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS film_graph_edges CASCADE;
DROP TABLE IF EXISTS film_graph_nodes CASCADE;
-- +goose StatementEnd
