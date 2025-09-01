## In-memory Key-Value store like Redis

### Possible commands:
- `get <key>`
    
    This returns the value (object with attributes and their values). Returns error if the key is not found.

- `put <key> <TTL> <attributeKey1> <attributeValue1> <attributeKey2> <attributeValue2>....`
    
    Adds tke key and the attributes to the key-value store. If the key already exists, then the value is replaced only if the data types of the attributes are maintained.
    Set TTL in seconds (do not pass 0), and -1 if no TTL is needed, meaning the key never gets evicted automatically.

- `delete <key>`

    Deletes the key. Returns nothing irrespective of whether the key is present or not.

- `search <attributeKey> <attributeValue>`

    Returns a list of all the keys that have the given attribute key-value pair.

- `keys`

    Returns a list of all the existing keys.

- `exit`

    Exits the code.

---

This is my solution to the machine coding question mentioned [here](https://workat.tech/machine-coding/practice/design-key-value-store-6gz6cq124k65).