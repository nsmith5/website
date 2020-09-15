---
title: "Using PGP Keys for SSH on Fedora"
date: 2020-09-14T22:31:21-07:00
draft: false
---

I learned recently that PGP keys can be used for SSH authentication. This is
quite convenient because if you're already managing a PGP key to sign git
commits (hint: you should probably be [doing this][signing-commits]), then
you'll have one less key to manage. As a bonus, if you've added your PGP public
key to [Github](https://github.com) you can immediately use a subkey with
authentication privileges to pull and push to your repositories over SSH.

[signing-commits]: https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work

This has been documented extensively else where, but the gist of the idea is
that `gpg-agent` can act as an OpenSSH compatibly SSH agent (the program used
to store and manage SSH keys). If you launch `gpg-agent` as follows:

```shell
$ eval $(gpg-agent --daemon --enable-ssh-support)
```

SSH will use your PGP keys to authenticate. Run `ssh-add -L` to list the
available keys to verify SSH found them. You should see keys of the form
`openpgp:{key-id}` or `cardno:{card-id}` if you're using a Smart Card like
[Yubikey](https://www.yubico.com/products/).

The only problem, is that on Gnome based systems there is already an SSH agent
running, namely [Gnome Keyring][gnome-keyring]. Most tutorials suggest
disabling this to some degree, but that isn't actually necessary. Moreover,
I've found that Gnome will just complain that the keyring hasn't started and
keeping asking me to authenticate to it when I've tried these methods. We just
need to make sure SSH finds the agent `gpg-agent` is running instead of the one
Gnome Keyring is running. How does SSH find its agent? The `SSH_AUTH_SOCK`
environment variable.

[gnome-keyring]: https://wiki.archlinux.org/index.php/GNOME/Keyring

Alright so here is the low down. To start, enable the `gpg-agent` systemd user
sockets to make sure the agent is enabled on start up.

```shell
$ systemctl --user enable --now gpg-agent.socket
$ sytsemctl --user enable --now gpg-agent-ssh.socket
```

Then, update `SSH_AUTH_SOCK` in your `bashrc` (or equivalent).

```shell
export SSH_AUTH_SOCK="$(gpgconf --list-dirs agent-ssh-socket)"
```

That's it! Now if you're lucky enough to have your PGP keys on a Yubikey, you
can completely remove all private keys from your laptop. A perfect example of
security in depth.

## More Reading

Wanna learn more? Here are some more articles:

- [Using PGP Keys for SSH Auth | Arch Wiki](https://wiki.archlinux.org/index.php/GnuPG#Using_a_PGP_key_for_SSH_authentication)
- [Gnome Keyring | Arch Wiki](https://wiki.archlinux.org/index.php/GNOME/Keyring)
- [Dr Duh's Excellent Yubikey Guide](https://github.com/drduh/YubiKey-Guide#ssh)
