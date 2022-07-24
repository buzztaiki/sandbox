const { Clutter } = imports.gi;

function log(message) {
  global.log(`[capture-event1]: ${message}`);
}

let handlerId = null;

function enable() {
  log("enabled");
  handlerId = global.stage.connect("captured-event", (_actor, event) => {
    if (event.type() == Clutter.EventType.BUTTON_PRESS) {
      log(`button=${event.get_button()}`);
    }
    return Clutter.EVENT_PROPAGATE;
  });
  log(`handlerId=${handlerId}`);
}

function disable() {
  log("disabled");
  global.stage.disconnect(handlerId);
  handlerId = null;
}
