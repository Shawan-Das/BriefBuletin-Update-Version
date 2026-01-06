-- SCHEMA: hrm
CREATE SCHEMA news;


CREATE TABLE news.users (
	user_id serial4 NOT NULL,
	user_name text NOT NULL,
	email text NOT NULL,
	phone text NOT NULL,
	pass text NOT NULL,
	pss_valid bool DEFAULT true NOT NULL,
	otp text NOT NULL,
	user_valid bool DEFAULT false NOT NULL,
	otp_exp timestamp NOT NULL,
	"role" text NOT NULL
);

CREATE TABLE news.categories (
	id serial4 NOT NULL,
	"name" text NOT NULL,
	slug text NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT categories_unique UNIQUE (name),
	CONSTRAINT categories_unique_1 UNIQUE (slug)
);

--    status ENUM('draft', 'published') DEFAULT 'draft',
CREATE TABLE news.articles (
	id serial4 NOT NULL,
	title text NOT NULL,
	summary text NOT NULL,
	"content" text NOT NULL,
	featured_image text NULL,
	category_id int4 NOT NULL,
	status text DEFAULT 'draft'::text NOT NULL,
	published_at timestamp NULL,
	views_count int8 DEFAULT 0 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	source_url text NOT NULL,
	CONSTRAINT articles_pk PRIMARY KEY (id),
	CONSTRAINT articles_categories_fk FOREIGN KEY (category_id) REFERENCES news.categories(id)
);

CREATE TABLE news."comments" (
	id serial4 NOT NULL,
	article_id int4 NOT NULL,
	user_name text NOT NULL,
	user_email text NOT NULL,
	"content" text NOT NULL,
	is_approved bool DEFAULT false NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT comments_unique UNIQUE (id),
	CONSTRAINT comments_articles_fk FOREIGN KEY (article_id) REFERENCES news.articles(id) ON DELETE CASCADE
);