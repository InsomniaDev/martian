import GoogleAssistant from "google-assistant";

// https://www.npmjs.com/package/google-assistant

class GoogleAssistant {
  constructor() {}

  startConversation() {
    const path = require("path");
    const GoogleAssistant = require("google-assistant");
    const config = {
      auth: {
        keyFilePath: path.resolve(__dirname, "YOUR_API_KEY_FILE_PATH.json"),
        // where you want the tokens to be saved
        // will create the directory if not already there
        savedTokensPath: path.resolve(__dirname, "tokens.json")
      },
      // this param is optional, but all options will be shown
      conversation: {
        audio: {
          encodingIn: "LINEAR16", // supported are LINEAR16 / FLAC (defaults to LINEAR16)
          sampleRateIn: 16000, // supported rates are between 16000-24000 (defaults to 16000)
          encodingOut: "LINEAR16", // supported are LINEAR16 / MP3 / OPUS_IN_OGG (defaults to LINEAR16)
          sampleRateOut: 24000 // supported are 16000 / 24000 (defaults to 24000)
        },
        lang: "en-US", // language code for input/output (defaults to en-US)
        deviceModelId: "xxxxxxxx", // use if you've gone through the Device Registration process
        deviceId: "xxxxxx", // use if you've gone through the Device Registration process
        deviceLocation: {
          coordinates: {
            // set the latitude and longitude of the device
            latitude: xxxxxx,
            longitude: xxxxx
          }
        },
        textQuery: "What time is it?", // if this is set, audio input is ignored
        isNew: true, // set this to true if you want to force a new conversation and ignore the old state
        screen: {
          isOn: true // set this to true if you want to output results to a screen
        }
      }
    };

    const assistant = new GoogleAssistant(config.auth);

    // starts a new conversation with the assistant
    const startConversation = conversation => {
      // setup the conversation and send data to it
      // for a full example, see `examples/mic-speaker.js`

      conversation
        .on("audio-data", data => {
          // do stuff with the audio data from the server
          // usually send it to some audio output / file
        })
        .on("end-of-utterance", () => {
          // do stuff when done speaking to the assistant
          // usually just stop your audio input
        })
        .on("transcription", data => {
          // do stuff with the words you are saying to the assistant
        })
        .on("response", text => {
          // do stuff with the text that the assistant said back
        })
        .on("volume-percent", percent => {
          // do stuff with a volume percent change (range from 1-100)
        })
        .on("device-action", action => {
          // if you've set this device up to handle actions, you'll get that here
        })
        .on("screen-data", screen => {
          // if the screen.isOn flag was set to true, you'll get the format and data of the output
        })
        .on("ended", (error, continueConversation) => {
          // once the conversation is ended, see if we need to follow up
          if (error) console.log("Conversation Ended Error:", error);
          else if (continueConversation) assistant.start();
          else console.log("Conversation Complete");
        })
        .on("error", error => console.error(error));
    };

    // will start a conversation and wait for audio data
    // as soon as it's ready
    assistant
      .on("ready", () => assistant.start())
      .on("started", startConversation);
  }

  sendMessage() {
    const assistant = new GoogleAssistant();

    //     assistant.message()
    //     global.config.assistants[k] = new GoogleAssistant(i)
    //     let assistant = global.config.assistants[k];
    //     assistant.on('ready', () => cb());
    //     assistant.on('error', (e) => {
    //       console.log(`❌ Assistant Error when activating user ${k}. Trying next user ❌ \n`);
    //       return cb();
    //     })
    //   }, function(err){
    //     if(err) return reject(err.message);
    //     (async() => {
    //       console.log(await terminalImage.file('./icon.png'))
    //       console.log(`Assistant Relay is now setup and running for${users.map(u => ` ${u}`)} \n`)
    //       console.log(`You can now visit ${global.config.baseUrl} in a browser, or send POST requests to it`);
    //     })();
  }
}
