// Run a docker container locally
// docker run -p 6379:6379 -it -d --rm -v /home/adam/Desktop/develop/rusty/data:/data redislabs/redisgraph:edge
// docker run -d --name redis -p 6379:6379 -v /home/adam/Desktop/develop/rusty/redisdata:/data --restart=unless-stopped redislabs/redisgraph --loadmodule /usr/lib/redis/modules/redisgraph.so --save 60 1
// docker run -v /home/adam/Desktop/develop/rusty/data:/db --network host redislabs/redisinsight

// q    `   1`

// Create the device and return the value
CREATE (a:Device { uuid: '', name: '', type: '', state: '' })
RETURN a

// Get the device
MATCH (a:Device { uuid: '' })
RETURN a

// Create the area name
CREATE (a:Area { name: '' })
RETURN a

// Get the area
MATCH (a:Area { name: '' })
RETURN a

// Get the devices associated to an area
MATCH (:Area) <-[:RESIDES_IN]- (d:Device)
RETURN d

// Create the relationship between the Area and Device
MATCH (a:Device { uuid: '' }), (b:Area { name: '' })
CREATE (a) -[r:RESIDES_IN]-> (b)


// EVENTS
// Create the event
CREATE (a:Event { uuid: '', time: '', state: '', device: 'uuid', occurrences: 1 })
RETURN a

// Match event to the device
MATCH (a:Event { uuid: '' }), (b:Device { uuid: ''})
CREATE (a) -[r:IS_DEVICE]-> (b)

// Match event to the next event
CREATE (a:Event { values }) -[:PRECEDES]-> (b:Event { values }) 

MERGE (a:Event { name: "test2" }) -[:PRECEDES]-> (b:Event { name: "test3" }) -[:PRECEDES]-> (c:Event { name: "test4" }) 
ON CREATE SET a.count=1 
ON MATCH SET a.count=a.count+1

MERGE (a:Event { name: "test3" }) -[:PRECEDES]-> (b:Event { name: "test5" }) 
ON CREATE SET a.count=1 
ON MATCH SET a.count=a.count+1

MERGE (a:Event { name: "test2" }) -[:PRECEDES]-> (b:Event { name: "test3" })
ON CREATE SET a.count=1 
ON MATCH SET a.count=a.count+1



MERGE (Time {name:"today"})
MATCH ,

ON CREATE SET b.count=1 
ON MATCH SET b.count=b.count+1

// Merge the first event with today
MERGE (a:Time{name:"today"}) 
MERGE (b:Event{ name: "test3" })
MERGE (a) <-[:OCCURED_AT]- (b)

// Create the relationship to all of the other events
MATCH (a:Event { name: "test3" })
MERGE (b:Event { name: "test4" })
MERGE (c:Event { name: "test5" })
MERGE (a) -[:PRECEDES]-> (b) -[:PRECEDES]-> (c)
ON CREATE SET a.count=1 
ON MATCH SET a.count=a.count+1
RETURN a

// Round two check
MERGE (a:Time{name:"today"}) 
MERGE (b:Event{ name: "test7" })
MERGE (a) <-[:OCCURED_AT]- (b)

// Create the relationship to all of the other events
MATCH (Time{name:"today"}) <-[:OCCURED_AT]- (a:Event { name: "test7" })
MERGE (b:Event { name: "test8" }) 
MERGE (c:Event { name: "test9" })
MERGE (a) -[:PRECEDES]-> (b) -[:PRECEDES]-> (c)
ON CREATE SET a.count=1 
ON MATCH SET a.count=a.count+1
RETURN a

MERGE (a:Time{name:"today"}) 
MERGE (b:Event{ name: "test4" })
MERGE (a) <-[:OCCURED_AT]- (b)

MATCH (Time{name:"today"}) <-[:OCCURED_AT]- (a:Event { name: "test4" })
MERGE (b:Event { name: "test7" })
MERGE (c:Event { name: "test9" })
MERGE (a) -[:PRECEDES]-> (b) -[:PRECEDES]-> (c)
ON CREATE SET a.count=1 
ON MATCH SET a.count=a.count+1
RETURN a

// Return all future events
MATCH (Time{name:"today"}) <-[:OCCURED_AT*1..50]- (a:Event) -[:PRECEDES*1..10]-> (b)
RETURN a,b

MATCH (Time{name:"today"}) <-[:OCCURED_AT*1..50]- (a:Event)
RETURN a
ORDER BY a.count DESC