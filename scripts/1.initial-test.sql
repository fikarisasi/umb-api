CREATE TABLE IF NOT EXISTS article (
  id SERIAL  NOT NULL,
  name varchar(200) NOT NULL,
  client varchar(200) NOT NULL,
  url varchar(200) NOT NULL,
  notes text NOT NULL,
  PRIMARY KEY (id)
);

insert into article (name, client, url, notes) values (
'qqq', 'ww', 'ee','rr');
