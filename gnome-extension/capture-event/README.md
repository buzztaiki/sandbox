# capture-event

Wayland だと Atspi が使えないから、 代わりに `global.stage` の `captured-event` を購読することで、全てのクリックをトラックしようとする試み。

https://github.com/harshadgavali/gnome-gesture-improvements が、以下のような事をやっていたので、実現できるんじゃないかと思った

```ts
this._stageCaptureEvent = global.stage.connect('captured-event::touchpad', this._handleEvent.bind(this));
```

が、結果として失敗。デスクトップのクリックはトラックできたのだけど、ウィンドウのクリックがトラックできない。つまり `global.stage` っていうのはそういう事。
全てのウィンドウに接続すればよいのだろうけど。むずい。

## References
- https://github.com/harshadgavali/gnome-gesture-improvements
- https://gjs-docs.gnome.org/clutter8~8_api/clutter.actor#signal-captured-event
- https://gjs-docs.gnome.org/clutter8~8_api/clutter.event
