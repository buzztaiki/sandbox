# flatpak-hello

https://docs.flatpak.org/en/latest/first-build.html


## build

```console
# add repository
% flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo

# build and install
% flatpak-builder --force-clean --user --install-deps-from=flathub --repo=repo --install builddir org.flatpak.Hello.yml
```

## run

```console
% flatpak run org.flatpak.Hello
Hello world, from a sandbox
```

## cleanup

```console
% flatpak remove org.flatpak.Hello
```
