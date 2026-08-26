CREATE TABLE IF NOT EXISTS graphs(id text primary key,name text not null,status text not null);
CREATE TABLE IF NOT EXISTS vertices(id text primary key,graph_id text not null,vertex_type text not null,properties jsonb not null,version bigint not null);
CREATE TABLE IF NOT EXISTS edges(id text primary key,graph_id text not null,from_id text not null,to_id text not null,edge_type text not null,properties jsonb not null,version bigint not null);
CREATE TABLE IF NOT EXISTS snapshots(id text primary key,graph_id text not null,commit_id bigint not null,created_at timestamptz not null);
