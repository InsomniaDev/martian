# martian
A smart home software that integrates with Hubitat and will learn your activity within your house and adjust accordingly.

## Skills
martian is currently capable of learning and setting automations that occur through:
- Actions
  - these automations are learned through a cause and effect sequence. _eg., if you walk into a room and the motion sensor goes active followed by a light being turned on, then there will be an automation triggered after enough occurrences to automatically turn on the light_
- Time
  - These automations occur when a light turns on repeatedly at a specific time. _eg., if you always turn on the kitchen light at 6am, then the light will automatically turn on at 6am_

martian is also capable of learning energy efficiency settings:
- By default each device is set to ten minutes energy efficiency. _eg., it will turn off after ten minutes_. You can exclude devices from energy efficiency by sending an API request to the running container to turn off energy efficiency with the device ID
- Auto-correction
  - If the light is turned off at ten minutes and immediately turned back on, then the energy efficiency interval will be increase
- Continuous tuning
  - The energy efficiency timeframe will be decreased daily to save energy

## Licensing
martian is an open source product licensed under GPLv3.