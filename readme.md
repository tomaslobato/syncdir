## how it works

client runs `sync`.
it lists every file and folder with it's content in json and sends it to server.
server iterates the json and runs mkdir or newfile for each of them, and completes the files with the contents received.
if the files already exist it's going to ignore them.
after that, it'll return the same thing, an infinite json of files and folders with content, which will then be populated into the client's folder.
and that would be it.

after every modification in the folder have to run:
`docker exec -u 33 -it nextcloud php occ files:scan --all`

for permissions of deleting and editing on webui
`sudo chown -R www-data:www-data ~/nextcloud`
`docker exec -it <container_id> chown -R www-data:www-data /var/www/html/data`