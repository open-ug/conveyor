---
sidebar_position: 1
---

# Resource Definitions

Conveyor CI stores all its information in objects called [Resources](resources). These objects are defined by the platform developer and can store data in any format, as long as it's predefined. In order to define this format, we use other internal objects called **Resource Definitions**.

Resource Definitions are objects that define the syntax used to create a [Resource](resources) object. When defining this syntax, Resource Definitions use the [JSON Schema (draft 4 / validation subset)](https://json-schema.org/specification-links#draft-4) as a schema language.

## Components of a Resource Definition

A Resource definition is composed of specific fields that all store certain information about a resource. These fields include:

- **Name**: The name field is used to specify the name of the Resource Definition. It is meant to be unique across different Resource Definitions. It is an alpha-numerical string field that does not include spaces or any special characters except “-”.
- **Description**: The description field defines the description of the Resource being defined. This can be information about what the Resource indicates or is intended to accomplish. It is a string that can include spaces
- **Version**: The version field indicates the version of the Resource Definition.
- **Schema**: The schema field includes the JSON Schema definition of the Resource.

## How Resource Definitions are used

When a Resource is sent to the API Server, the API Server will use the Resource Definition to validate the Resource syntax. It uses a JSON Schema validator and returns a detailed error in case wrong syntax is detected. Otherwise the Resource is saved in Conveyor CI.

## Creating a Resource Definition

When creating a Resource Definition. The first step is to research what data the Resource is expected to hold. Then a suitable JSON format that can handle the data is created. Next a JSON Schema is defined for that format and placed in the `schema` field of the Resource Definition object.
Take an example of this Resource Definition

```json
{
  "name": "workflow",
  "version": "0.0.1",
  "schema": {
    "pipeline": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "stages": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "distributed": {
          "type": "boolean"
        },
        "runners": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      },
      "required": [
        "name",
        "stages",
        "distributed",
        "runners"
      ]
    }
  }
}
```

Lastly the resource is sent to the API Server using a POST request to the `/resource-definitions/` route.
