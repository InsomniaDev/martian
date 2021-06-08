import { EventsGraph } from "./knowledge/modules/event";

// Example of event
// {
//     device: "lutron",
//     time: new Date(),
//     state: "off",
//     id: "1"
// }

const eventData = [];
export function AlertTheBrain(data) {
  data.time = new Date();
  // Will remove the event if they are the same event
  for (let i = 0; i < eventData.length; i++) {
    if (eventData[i].id === data.id) {
      eventData.splice(i, 1);
      checkEvents();
      return;
    }
  }
  eventData.push(data);
  checkEvents();
}

// FIXME: Should we put all of this logic into the events itself?
// Make a random id to insert into the event database
function makeid(length) {
  var result = "";
  var characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

// Cycle through the events as they occur and the server is active
function checkEvents() {
  // Check if there are more than one event and that there are only two events
  if (eventData.length > 0 && eventData.length === 2) {
    // Go through the logic to parse the following events

    // Get the time difference in seconds
    const timeDifference = (eventData[1].time - eventData[0].time) / 1000;
    // If the events are more than 30 seconds apart then don't record the event
    if (timeDifference > 30) {
      eventData.shift();
      return;
    }

    // Check if there are more than one events recorded for the last minute
    if (eventData.length === 1) {
      // Remove the single event then jump out of the processing logic
      eventData.pop();
    }

    // Get the first event and start creating the cypher statements
    const { device, state, id } = eventData.shift();
    const { device: nextDevice, state: nextState, id: nextId } = eventData[0];
    const cipherEventForCreation = `
      MERGE (a:event{device:'${device}', state:'${state}', deviceId:'${id}'})
      MERGE (b:event{device:'${nextDevice}', state:'${nextState}', deviceId:'${nextId}'})
      MERGE (a) -[c:RELATES]-> (b)
      ON CREATE SET c.weight=1 
      ON MATCH SET c.weight=c.weight+1
    `;

    // Insert into the database the events that are occurring
    EventsGraph.InsertIntoGraph(cipherEventForCreation);
  }
}
