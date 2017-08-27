#Deploying with ansible
To install it on your server with ansible, first replace on you machine the /etc/ansible/hosts file with the hosts file in this folder, and place your server IP there on the hosts file. Then, you need to replace the `image_source_url` property inside the melanite role [defaults file](https://github.com/jademcosta/melanite/blob/develop/deploy/ansible/roles/melanite/defaults/main.yml).

After, run `ansible-playbook ansible-playbook deploy/ansible/playbook.yml`. That's it, Melanite should be running on your server.
If it refuses to start, you can check systemd logs with: `journalctl -u melanite`.

This ansible playbook was tested on Ubuntu 16.04, but should work OK on new Ubuntu versions.
