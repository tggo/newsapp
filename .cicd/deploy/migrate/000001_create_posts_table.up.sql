create table jrnl_posts
(
    id         serial
        constraint jrnl_posts_pk
            primary key,
    created_at timestamp without time zone default now() not null,
    updated_at timestamp without time zone,
    title      text                                      not null,
    body       text                                      not null
);

comment on table jrnl_posts is 'news posts';

create unique index jrnl_posts_id_uindex
    on jrnl_posts (id);

insert into jrnl_posts(title, body) VALUES(
    'Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit...',
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis ultrices enim id leo eleifend, ' ||
    'commodo luctus lorem venenatis. Morbi et mi vitae sapien condimentum egestas ac ac felis. ' ||
    'Ut euismod justo sapien. Phasellus at quam quis mauris mollis varius. Nulla arcu neque, ornare id ex eget, ' ||
    'rhoncus mattis purus. Phasellus id lectus odio. Curabitur blandit tempus purus ac finibus. ' ||
    'Pellentesque vulputate erat nec nisl sagittis, eu varius nulla tincidunt. Duis laoreet egestas massa sit' ||
    ' amet fermentum. Morbi egestas auctor tincidunt. Nam tincidunt lectus ut augue rutrum convallis. ' ||
    'Nunc sed tellus iaculis, aliquam mauris ut, eleifend velit.')
