import { app, EventGridEvent, InvocationContext } from "@azure/functions";

export async function EventGridExample(event: EventGridEvent, context: InvocationContext): Promise<void> {
    context.log('Event grid function processed event:', event);
}

app.eventGrid('EventGridExample', {
    handler: EventGridExample
});
