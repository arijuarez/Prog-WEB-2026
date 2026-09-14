-- Created by Redgate Data Modeler (https://datamodeler.redgate-platform.com)
-- Last modification date: 2026-09-12 22:03:08.016

-- tables
-- Table: Album
CREATE TABLE Album (
    id_album serial  NOT NULL,
    album_name varchar  NOT NULL,
    album_dur int  NOT NULL,
    album_ph int  NOT NULL,
    release_date date  NOT NULL,
    Artist_id_artist int  NOT NULL,
    LIkes int NOT NULL DEFAULT 0,
    CONSTRAINT Album_pk PRIMARY KEY (id_album)
);

-- Table: Artist
CREATE TABLE Artist (
    id_artist serial  NOT NULL,
    artist_name varchar  NOT NULL,
    photo int  NOT NULL,
    artist_bio text  NOT NULL,
    CONSTRAINT Artist_pk PRIMARY KEY (id_artist)
);

-- Table: Artist_Song
CREATE TABLE Artist_Song (
    Artist_id_artist int  NOT NULL,
    Song_id_song int  NOT NULL,
    Roll bool  NOT NULL,
    CONSTRAINT Artist_Song_pk PRIMARY KEY (Artist_id_artist,Song_id_song)
);

-- Table: Review
CREATE TABLE Review (
    id_review serial  NOT NULL,
    review_text text  NOT NULL,
    review_rate float  NOT NULL,
    review_date date  NOT NULL,
    User_id_user int  NOT NULL,
    Song_id_song int  NULL,
    Album_id_album int  NULL,
    Likes int NOT NULL DEFAULT 0,
    CONSTRAINT Review_pk PRIMARY KEY (id_review)
);

-- Table: Song
CREATE TABLE Song (
    id_song serial  NOT NULL,
    song_name varchar  NOT NULL,
    song_dur int  NOT NULL,
    likes int NOT NULL DEFAULT 0,
    Album_id_album int  NOT NULL,
    CONSTRAINT Song_pk PRIMARY KEY (id_song)
);

-- Table: User
CREATE TABLE "User" (
    id_user serial  NOT NULL,
    user_name varchar  NOT NULL,
    amigos int  NOT NULL,
    photo int  NOT NULL,
    CONSTRAINT User_pk PRIMARY KEY (id_user)
);

-- foreign keys
-- Reference: Album_Artist (table: Album)
ALTER TABLE Album ADD CONSTRAINT Album_Artist
    FOREIGN KEY (Artist_id_artist)
    REFERENCES Artist (id_artist)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Artist_Song_Artist (table: Artist_Song)
ALTER TABLE Artist_Song ADD CONSTRAINT Artist_Song_Artist
    FOREIGN KEY (Artist_id_artist)
    REFERENCES Artist (id_artist)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Artist_Song_Song (table: Artist_Song)
ALTER TABLE Artist_Song ADD CONSTRAINT Artist_Song_Song
    FOREIGN KEY (Song_id_song)
    REFERENCES Song (id_song)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Review_Album (table: Review)
ALTER TABLE Review ADD CONSTRAINT Review_Album
    FOREIGN KEY (Album_id_album)
    REFERENCES Album (id_album)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Review_Song (table: Review)
ALTER TABLE Review ADD CONSTRAINT Review_Song
    FOREIGN KEY (Song_id_song)
    REFERENCES Song (id_song)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Review_User (table: Review)
ALTER TABLE Review ADD CONSTRAINT Review_User
    FOREIGN KEY (User_id_user)
    REFERENCES "User" (id_user)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Song_Album (table: Song)
ALTER TABLE Song ADD CONSTRAINT Song_Album
    FOREIGN KEY (Album_id_album)
    REFERENCES Album (id_album)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- End of file.