# Generated TypeScript README
This README will guide you through the process of using the generated JavaScript SDK package for the connector `driftlock`. It will also provide examples on how to use your generated SDK to call your Data Connect queries and mutations.

***NOTE:** This README is generated alongside the generated SDK. If you make changes to this file, they will be overwritten when the SDK is regenerated.*

# Table of Contents
- [**Overview**](#generated-javascript-readme)
- [**Accessing the connector**](#accessing-the-connector)
  - [*Connecting to the local Emulator*](#connecting-to-the-local-emulator)
- [**Queries**](#queries)
  - [*GetUser*](#getuser)
  - [*ListUsersByEmail*](#listusersbyemail)
  - [*GetDataset*](#getdataset)
  - [*ListDatasetsByUser*](#listdatasetsbyuser)
  - [*GetModelConfiguration*](#getmodelconfiguration)
  - [*ListModelConfigurationsByUser*](#listmodelconfigurationsbyuser)
  - [*GetDetectionTask*](#getdetectiontask)
  - [*ListDetectionTasksByUser*](#listdetectiontasksbyuser)
  - [*GetAnomaliesByTask*](#getanomaliesbytask)
  - [*GetHighScoreAnomalies*](#gethighscoreanomalies)
- [**Mutations**](#mutations)
  - [*CreateUser*](#createuser)
  - [*UpdateUser*](#updateuser)
  - [*CreateDataset*](#createdataset)
  - [*UpdateDatasetStatus*](#updatedatasetstatus)
  - [*UpdateDataset*](#updatedataset)
  - [*CreateModelConfiguration*](#createmodelconfiguration)
  - [*UpdateModelConfiguration*](#updatemodelconfiguration)
  - [*CreateDetectionTask*](#createdetectiontask)
  - [*UpdateDetectionTask*](#updatedetectiontask)
  - [*CompleteDetectionTask*](#completedetectiontask)
  - [*CreateAnomaly*](#createanomaly)

# Accessing the connector
A connector is a collection of Queries and Mutations. One SDK is generated for each connector - this SDK is generated for the connector `driftlock`. You can find more information about connectors in the [Data Connect documentation](https://firebase.google.com/docs/data-connect#how-does).

You can use this generated SDK by importing from the package `@driftlock/dataconnect` as shown below. Both CommonJS and ESM imports are supported.

You can also follow the instructions from the [Data Connect documentation](https://firebase.google.com/docs/data-connect/web-sdk#set-client).

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig } from '@driftlock/dataconnect';

const dataConnect = getDataConnect(connectorConfig);
```

## Connecting to the local Emulator
By default, the connector will connect to the production service.

To connect to the emulator, you can use the following code.
You can also follow the emulator instructions from the [Data Connect documentation](https://firebase.google.com/docs/data-connect/web-sdk#instrument-clients).

```typescript
import { connectDataConnectEmulator, getDataConnect } from 'firebase/data-connect';
import { connectorConfig } from '@driftlock/dataconnect';

const dataConnect = getDataConnect(connectorConfig);
connectDataConnectEmulator(dataConnect, 'localhost', 9399);
```

After it's initialized, you can call your Data Connect [queries](#queries) and [mutations](#mutations) from your generated SDK.

# Queries

There are two ways to execute a Data Connect Query using the generated Web SDK:
- Using a Query Reference function, which returns a `QueryRef`
  - The `QueryRef` can be used as an argument to `executeQuery()`, which will execute the Query and return a `QueryPromise`
- Using an action shortcut function, which returns a `QueryPromise`
  - Calling the action shortcut function will execute the Query and return a `QueryPromise`

The following is true for both the action shortcut function and the `QueryRef` function:
- The `QueryPromise` returned will resolve to the result of the Query once it has finished executing
- If the Query accepts arguments, both the action shortcut function and the `QueryRef` function accept a single argument: an object that contains all the required variables (and the optional variables) for the Query
- Both functions can be called with or without passing in a `DataConnect` instance as an argument. If no `DataConnect` argument is passed in, then the generated SDK will call `getDataConnect(connectorConfig)` behind the scenes for you.

Below are examples of how to use the `driftlock` connector's generated functions to execute each query. You can also follow the examples from the [Data Connect documentation](https://firebase.google.com/docs/data-connect/web-sdk#using-queries).

## GetUser
You can execute the `GetUser` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
getUser(vars: GetUserVariables): QueryPromise<GetUserData, GetUserVariables>;

interface GetUserRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetUserVariables): QueryRef<GetUserData, GetUserVariables>;
}
export const getUserRef: GetUserRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
getUser(dc: DataConnect, vars: GetUserVariables): QueryPromise<GetUserData, GetUserVariables>;

interface GetUserRef {
  ...
  (dc: DataConnect, vars: GetUserVariables): QueryRef<GetUserData, GetUserVariables>;
}
export const getUserRef: GetUserRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the getUserRef:
```typescript
const name = getUserRef.operationName;
console.log(name);
```

### Variables
The `GetUser` query requires an argument of type `GetUserVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface GetUserVariables {
  id: UUIDString;
}
```
### Return Type
Recall that executing the `GetUser` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `GetUserData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface GetUserData {
  user?: {
    id: UUIDString;
    displayName: string;
    email?: string | null;
    photoUrl?: string | null;
    createdAt: TimestampString;
  } & User_Key;
}
```
### Using `GetUser`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, getUser, GetUserVariables } from '@driftlock/dataconnect';

// The `GetUser` query requires an argument of type `GetUserVariables`:
const getUserVars: GetUserVariables = {
  id: ..., 
};

// Call the `getUser()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await getUser(getUserVars);
// Variables can be defined inline as well.
const { data } = await getUser({ id: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await getUser(dataConnect, getUserVars);

console.log(data.user);

// Or, you can use the `Promise` API.
getUser(getUserVars).then((response) => {
  const data = response.data;
  console.log(data.user);
});
```

### Using `GetUser`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, getUserRef, GetUserVariables } from '@driftlock/dataconnect';

// The `GetUser` query requires an argument of type `GetUserVariables`:
const getUserVars: GetUserVariables = {
  id: ..., 
};

// Call the `getUserRef()` function to get a reference to the query.
const ref = getUserRef(getUserVars);
// Variables can be defined inline as well.
const ref = getUserRef({ id: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = getUserRef(dataConnect, getUserVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.user);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.user);
});
```

## ListUsersByEmail
You can execute the `ListUsersByEmail` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
listUsersByEmail(vars: ListUsersByEmailVariables): QueryPromise<ListUsersByEmailData, ListUsersByEmailVariables>;

interface ListUsersByEmailRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListUsersByEmailVariables): QueryRef<ListUsersByEmailData, ListUsersByEmailVariables>;
}
export const listUsersByEmailRef: ListUsersByEmailRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
listUsersByEmail(dc: DataConnect, vars: ListUsersByEmailVariables): QueryPromise<ListUsersByEmailData, ListUsersByEmailVariables>;

interface ListUsersByEmailRef {
  ...
  (dc: DataConnect, vars: ListUsersByEmailVariables): QueryRef<ListUsersByEmailData, ListUsersByEmailVariables>;
}
export const listUsersByEmailRef: ListUsersByEmailRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the listUsersByEmailRef:
```typescript
const name = listUsersByEmailRef.operationName;
console.log(name);
```

### Variables
The `ListUsersByEmail` query requires an argument of type `ListUsersByEmailVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface ListUsersByEmailVariables {
  email: string;
}
```
### Return Type
Recall that executing the `ListUsersByEmail` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `ListUsersByEmailData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface ListUsersByEmailData {
  users: ({
    id: UUIDString;
    displayName: string;
    email?: string | null;
    photoUrl?: string | null;
    createdAt: TimestampString;
  } & User_Key)[];
}
```
### Using `ListUsersByEmail`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, listUsersByEmail, ListUsersByEmailVariables } from '@driftlock/dataconnect';

// The `ListUsersByEmail` query requires an argument of type `ListUsersByEmailVariables`:
const listUsersByEmailVars: ListUsersByEmailVariables = {
  email: ..., 
};

// Call the `listUsersByEmail()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await listUsersByEmail(listUsersByEmailVars);
// Variables can be defined inline as well.
const { data } = await listUsersByEmail({ email: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await listUsersByEmail(dataConnect, listUsersByEmailVars);

console.log(data.users);

// Or, you can use the `Promise` API.
listUsersByEmail(listUsersByEmailVars).then((response) => {
  const data = response.data;
  console.log(data.users);
});
```

### Using `ListUsersByEmail`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, listUsersByEmailRef, ListUsersByEmailVariables } from '@driftlock/dataconnect';

// The `ListUsersByEmail` query requires an argument of type `ListUsersByEmailVariables`:
const listUsersByEmailVars: ListUsersByEmailVariables = {
  email: ..., 
};

// Call the `listUsersByEmailRef()` function to get a reference to the query.
const ref = listUsersByEmailRef(listUsersByEmailVars);
// Variables can be defined inline as well.
const ref = listUsersByEmailRef({ email: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = listUsersByEmailRef(dataConnect, listUsersByEmailVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.users);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.users);
});
```

## GetDataset
You can execute the `GetDataset` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
getDataset(vars: GetDatasetVariables): QueryPromise<GetDatasetData, GetDatasetVariables>;

interface GetDatasetRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetDatasetVariables): QueryRef<GetDatasetData, GetDatasetVariables>;
}
export const getDatasetRef: GetDatasetRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
getDataset(dc: DataConnect, vars: GetDatasetVariables): QueryPromise<GetDatasetData, GetDatasetVariables>;

interface GetDatasetRef {
  ...
  (dc: DataConnect, vars: GetDatasetVariables): QueryRef<GetDatasetData, GetDatasetVariables>;
}
export const getDatasetRef: GetDatasetRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the getDatasetRef:
```typescript
const name = getDatasetRef.operationName;
console.log(name);
```

### Variables
The `GetDataset` query requires an argument of type `GetDatasetVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface GetDatasetVariables {
  id: UUIDString;
}
```
### Return Type
Recall that executing the `GetDataset` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `GetDatasetData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface GetDatasetData {
  dataset?: {
    id: UUIDString;
    name: string;
    filename: string;
    uploadDate: TimestampString;
    fileLocation: string;
    status: string;
    description?: string | null;
    columnInfo?: string | null;
    user: {
      id: UUIDString;
      displayName: string;
      email?: string | null;
    } & User_Key;
  } & Dataset_Key;
}
```
### Using `GetDataset`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, getDataset, GetDatasetVariables } from '@driftlock/dataconnect';

// The `GetDataset` query requires an argument of type `GetDatasetVariables`:
const getDatasetVars: GetDatasetVariables = {
  id: ..., 
};

// Call the `getDataset()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await getDataset(getDatasetVars);
// Variables can be defined inline as well.
const { data } = await getDataset({ id: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await getDataset(dataConnect, getDatasetVars);

console.log(data.dataset);

// Or, you can use the `Promise` API.
getDataset(getDatasetVars).then((response) => {
  const data = response.data;
  console.log(data.dataset);
});
```

### Using `GetDataset`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, getDatasetRef, GetDatasetVariables } from '@driftlock/dataconnect';

// The `GetDataset` query requires an argument of type `GetDatasetVariables`:
const getDatasetVars: GetDatasetVariables = {
  id: ..., 
};

// Call the `getDatasetRef()` function to get a reference to the query.
const ref = getDatasetRef(getDatasetVars);
// Variables can be defined inline as well.
const ref = getDatasetRef({ id: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = getDatasetRef(dataConnect, getDatasetVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.dataset);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.dataset);
});
```

## ListDatasetsByUser
You can execute the `ListDatasetsByUser` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
listDatasetsByUser(vars: ListDatasetsByUserVariables): QueryPromise<ListDatasetsByUserData, ListDatasetsByUserVariables>;

interface ListDatasetsByUserRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListDatasetsByUserVariables): QueryRef<ListDatasetsByUserData, ListDatasetsByUserVariables>;
}
export const listDatasetsByUserRef: ListDatasetsByUserRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
listDatasetsByUser(dc: DataConnect, vars: ListDatasetsByUserVariables): QueryPromise<ListDatasetsByUserData, ListDatasetsByUserVariables>;

interface ListDatasetsByUserRef {
  ...
  (dc: DataConnect, vars: ListDatasetsByUserVariables): QueryRef<ListDatasetsByUserData, ListDatasetsByUserVariables>;
}
export const listDatasetsByUserRef: ListDatasetsByUserRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the listDatasetsByUserRef:
```typescript
const name = listDatasetsByUserRef.operationName;
console.log(name);
```

### Variables
The `ListDatasetsByUser` query requires an argument of type `ListDatasetsByUserVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface ListDatasetsByUserVariables {
  userId: UUIDString;
}
```
### Return Type
Recall that executing the `ListDatasetsByUser` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `ListDatasetsByUserData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface ListDatasetsByUserData {
  datasets: ({
    id: UUIDString;
    name: string;
    filename: string;
    uploadDate: TimestampString;
    fileLocation: string;
    status: string;
    description?: string | null;
    columnInfo?: string | null;
  } & Dataset_Key)[];
}
```
### Using `ListDatasetsByUser`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, listDatasetsByUser, ListDatasetsByUserVariables } from '@driftlock/dataconnect';

// The `ListDatasetsByUser` query requires an argument of type `ListDatasetsByUserVariables`:
const listDatasetsByUserVars: ListDatasetsByUserVariables = {
  userId: ..., 
};

// Call the `listDatasetsByUser()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await listDatasetsByUser(listDatasetsByUserVars);
// Variables can be defined inline as well.
const { data } = await listDatasetsByUser({ userId: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await listDatasetsByUser(dataConnect, listDatasetsByUserVars);

console.log(data.datasets);

// Or, you can use the `Promise` API.
listDatasetsByUser(listDatasetsByUserVars).then((response) => {
  const data = response.data;
  console.log(data.datasets);
});
```

### Using `ListDatasetsByUser`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, listDatasetsByUserRef, ListDatasetsByUserVariables } from '@driftlock/dataconnect';

// The `ListDatasetsByUser` query requires an argument of type `ListDatasetsByUserVariables`:
const listDatasetsByUserVars: ListDatasetsByUserVariables = {
  userId: ..., 
};

// Call the `listDatasetsByUserRef()` function to get a reference to the query.
const ref = listDatasetsByUserRef(listDatasetsByUserVars);
// Variables can be defined inline as well.
const ref = listDatasetsByUserRef({ userId: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = listDatasetsByUserRef(dataConnect, listDatasetsByUserVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.datasets);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.datasets);
});
```

## GetModelConfiguration
You can execute the `GetModelConfiguration` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
getModelConfiguration(vars: GetModelConfigurationVariables): QueryPromise<GetModelConfigurationData, GetModelConfigurationVariables>;

interface GetModelConfigurationRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetModelConfigurationVariables): QueryRef<GetModelConfigurationData, GetModelConfigurationVariables>;
}
export const getModelConfigurationRef: GetModelConfigurationRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
getModelConfiguration(dc: DataConnect, vars: GetModelConfigurationVariables): QueryPromise<GetModelConfigurationData, GetModelConfigurationVariables>;

interface GetModelConfigurationRef {
  ...
  (dc: DataConnect, vars: GetModelConfigurationVariables): QueryRef<GetModelConfigurationData, GetModelConfigurationVariables>;
}
export const getModelConfigurationRef: GetModelConfigurationRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the getModelConfigurationRef:
```typescript
const name = getModelConfigurationRef.operationName;
console.log(name);
```

### Variables
The `GetModelConfiguration` query requires an argument of type `GetModelConfigurationVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface GetModelConfigurationVariables {
  id: UUIDString;
}
```
### Return Type
Recall that executing the `GetModelConfiguration` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `GetModelConfigurationData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface GetModelConfigurationData {
  modelConfiguration?: {
    id: UUIDString;
    name: string;
    algorithmType: string;
    createdAt: TimestampString;
    parameters?: string | null;
    description?: string | null;
    user: {
      id: UUIDString;
      displayName: string;
    } & User_Key;
  } & ModelConfiguration_Key;
}
```
### Using `GetModelConfiguration`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, getModelConfiguration, GetModelConfigurationVariables } from '@driftlock/dataconnect';

// The `GetModelConfiguration` query requires an argument of type `GetModelConfigurationVariables`:
const getModelConfigurationVars: GetModelConfigurationVariables = {
  id: ..., 
};

// Call the `getModelConfiguration()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await getModelConfiguration(getModelConfigurationVars);
// Variables can be defined inline as well.
const { data } = await getModelConfiguration({ id: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await getModelConfiguration(dataConnect, getModelConfigurationVars);

console.log(data.modelConfiguration);

// Or, you can use the `Promise` API.
getModelConfiguration(getModelConfigurationVars).then((response) => {
  const data = response.data;
  console.log(data.modelConfiguration);
});
```

### Using `GetModelConfiguration`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, getModelConfigurationRef, GetModelConfigurationVariables } from '@driftlock/dataconnect';

// The `GetModelConfiguration` query requires an argument of type `GetModelConfigurationVariables`:
const getModelConfigurationVars: GetModelConfigurationVariables = {
  id: ..., 
};

// Call the `getModelConfigurationRef()` function to get a reference to the query.
const ref = getModelConfigurationRef(getModelConfigurationVars);
// Variables can be defined inline as well.
const ref = getModelConfigurationRef({ id: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = getModelConfigurationRef(dataConnect, getModelConfigurationVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.modelConfiguration);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.modelConfiguration);
});
```

## ListModelConfigurationsByUser
You can execute the `ListModelConfigurationsByUser` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
listModelConfigurationsByUser(vars: ListModelConfigurationsByUserVariables): QueryPromise<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;

interface ListModelConfigurationsByUserRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListModelConfigurationsByUserVariables): QueryRef<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;
}
export const listModelConfigurationsByUserRef: ListModelConfigurationsByUserRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
listModelConfigurationsByUser(dc: DataConnect, vars: ListModelConfigurationsByUserVariables): QueryPromise<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;

interface ListModelConfigurationsByUserRef {
  ...
  (dc: DataConnect, vars: ListModelConfigurationsByUserVariables): QueryRef<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;
}
export const listModelConfigurationsByUserRef: ListModelConfigurationsByUserRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the listModelConfigurationsByUserRef:
```typescript
const name = listModelConfigurationsByUserRef.operationName;
console.log(name);
```

### Variables
The `ListModelConfigurationsByUser` query requires an argument of type `ListModelConfigurationsByUserVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface ListModelConfigurationsByUserVariables {
  userId: UUIDString;
}
```
### Return Type
Recall that executing the `ListModelConfigurationsByUser` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `ListModelConfigurationsByUserData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface ListModelConfigurationsByUserData {
  modelConfigurations: ({
    id: UUIDString;
    name: string;
    algorithmType: string;
    createdAt: TimestampString;
    parameters?: string | null;
    description?: string | null;
  } & ModelConfiguration_Key)[];
}
```
### Using `ListModelConfigurationsByUser`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, listModelConfigurationsByUser, ListModelConfigurationsByUserVariables } from '@driftlock/dataconnect';

// The `ListModelConfigurationsByUser` query requires an argument of type `ListModelConfigurationsByUserVariables`:
const listModelConfigurationsByUserVars: ListModelConfigurationsByUserVariables = {
  userId: ..., 
};

// Call the `listModelConfigurationsByUser()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await listModelConfigurationsByUser(listModelConfigurationsByUserVars);
// Variables can be defined inline as well.
const { data } = await listModelConfigurationsByUser({ userId: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await listModelConfigurationsByUser(dataConnect, listModelConfigurationsByUserVars);

console.log(data.modelConfigurations);

// Or, you can use the `Promise` API.
listModelConfigurationsByUser(listModelConfigurationsByUserVars).then((response) => {
  const data = response.data;
  console.log(data.modelConfigurations);
});
```

### Using `ListModelConfigurationsByUser`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, listModelConfigurationsByUserRef, ListModelConfigurationsByUserVariables } from '@driftlock/dataconnect';

// The `ListModelConfigurationsByUser` query requires an argument of type `ListModelConfigurationsByUserVariables`:
const listModelConfigurationsByUserVars: ListModelConfigurationsByUserVariables = {
  userId: ..., 
};

// Call the `listModelConfigurationsByUserRef()` function to get a reference to the query.
const ref = listModelConfigurationsByUserRef(listModelConfigurationsByUserVars);
// Variables can be defined inline as well.
const ref = listModelConfigurationsByUserRef({ userId: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = listModelConfigurationsByUserRef(dataConnect, listModelConfigurationsByUserVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.modelConfigurations);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.modelConfigurations);
});
```

## GetDetectionTask
You can execute the `GetDetectionTask` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
getDetectionTask(vars: GetDetectionTaskVariables): QueryPromise<GetDetectionTaskData, GetDetectionTaskVariables>;

interface GetDetectionTaskRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetDetectionTaskVariables): QueryRef<GetDetectionTaskData, GetDetectionTaskVariables>;
}
export const getDetectionTaskRef: GetDetectionTaskRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
getDetectionTask(dc: DataConnect, vars: GetDetectionTaskVariables): QueryPromise<GetDetectionTaskData, GetDetectionTaskVariables>;

interface GetDetectionTaskRef {
  ...
  (dc: DataConnect, vars: GetDetectionTaskVariables): QueryRef<GetDetectionTaskData, GetDetectionTaskVariables>;
}
export const getDetectionTaskRef: GetDetectionTaskRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the getDetectionTaskRef:
```typescript
const name = getDetectionTaskRef.operationName;
console.log(name);
```

### Variables
The `GetDetectionTask` query requires an argument of type `GetDetectionTaskVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface GetDetectionTaskVariables {
  id: UUIDString;
}
```
### Return Type
Recall that executing the `GetDetectionTask` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `GetDetectionTaskData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface GetDetectionTaskData {
  detectionTask?: {
    id: UUIDString;
    taskName: string;
    startDate: TimestampString;
    endDate?: TimestampString | null;
    status: string;
    resultsLocation?: string | null;
    notes?: string | null;
    user: {
      id: UUIDString;
      displayName: string;
    } & User_Key;
      dataset: {
        id: UUIDString;
        name: string;
        filename: string;
      } & Dataset_Key;
        modelConfiguration: {
          id: UUIDString;
          name: string;
          algorithmType: string;
        } & ModelConfiguration_Key;
  } & DetectionTask_Key;
}
```
### Using `GetDetectionTask`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, getDetectionTask, GetDetectionTaskVariables } from '@driftlock/dataconnect';

// The `GetDetectionTask` query requires an argument of type `GetDetectionTaskVariables`:
const getDetectionTaskVars: GetDetectionTaskVariables = {
  id: ..., 
};

// Call the `getDetectionTask()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await getDetectionTask(getDetectionTaskVars);
// Variables can be defined inline as well.
const { data } = await getDetectionTask({ id: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await getDetectionTask(dataConnect, getDetectionTaskVars);

console.log(data.detectionTask);

// Or, you can use the `Promise` API.
getDetectionTask(getDetectionTaskVars).then((response) => {
  const data = response.data;
  console.log(data.detectionTask);
});
```

### Using `GetDetectionTask`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, getDetectionTaskRef, GetDetectionTaskVariables } from '@driftlock/dataconnect';

// The `GetDetectionTask` query requires an argument of type `GetDetectionTaskVariables`:
const getDetectionTaskVars: GetDetectionTaskVariables = {
  id: ..., 
};

// Call the `getDetectionTaskRef()` function to get a reference to the query.
const ref = getDetectionTaskRef(getDetectionTaskVars);
// Variables can be defined inline as well.
const ref = getDetectionTaskRef({ id: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = getDetectionTaskRef(dataConnect, getDetectionTaskVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.detectionTask);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.detectionTask);
});
```

## ListDetectionTasksByUser
You can execute the `ListDetectionTasksByUser` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
listDetectionTasksByUser(vars: ListDetectionTasksByUserVariables): QueryPromise<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;

interface ListDetectionTasksByUserRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListDetectionTasksByUserVariables): QueryRef<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;
}
export const listDetectionTasksByUserRef: ListDetectionTasksByUserRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
listDetectionTasksByUser(dc: DataConnect, vars: ListDetectionTasksByUserVariables): QueryPromise<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;

interface ListDetectionTasksByUserRef {
  ...
  (dc: DataConnect, vars: ListDetectionTasksByUserVariables): QueryRef<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;
}
export const listDetectionTasksByUserRef: ListDetectionTasksByUserRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the listDetectionTasksByUserRef:
```typescript
const name = listDetectionTasksByUserRef.operationName;
console.log(name);
```

### Variables
The `ListDetectionTasksByUser` query requires an argument of type `ListDetectionTasksByUserVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface ListDetectionTasksByUserVariables {
  userId: UUIDString;
}
```
### Return Type
Recall that executing the `ListDetectionTasksByUser` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `ListDetectionTasksByUserData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface ListDetectionTasksByUserData {
  detectionTasks: ({
    id: UUIDString;
    taskName: string;
    startDate: TimestampString;
    endDate?: TimestampString | null;
    status: string;
    resultsLocation?: string | null;
    notes?: string | null;
    dataset: {
      id: UUIDString;
      name: string;
    } & Dataset_Key;
      modelConfiguration: {
        id: UUIDString;
        name: string;
      } & ModelConfiguration_Key;
  } & DetectionTask_Key)[];
}
```
### Using `ListDetectionTasksByUser`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, listDetectionTasksByUser, ListDetectionTasksByUserVariables } from '@driftlock/dataconnect';

// The `ListDetectionTasksByUser` query requires an argument of type `ListDetectionTasksByUserVariables`:
const listDetectionTasksByUserVars: ListDetectionTasksByUserVariables = {
  userId: ..., 
};

// Call the `listDetectionTasksByUser()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await listDetectionTasksByUser(listDetectionTasksByUserVars);
// Variables can be defined inline as well.
const { data } = await listDetectionTasksByUser({ userId: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await listDetectionTasksByUser(dataConnect, listDetectionTasksByUserVars);

console.log(data.detectionTasks);

// Or, you can use the `Promise` API.
listDetectionTasksByUser(listDetectionTasksByUserVars).then((response) => {
  const data = response.data;
  console.log(data.detectionTasks);
});
```

### Using `ListDetectionTasksByUser`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, listDetectionTasksByUserRef, ListDetectionTasksByUserVariables } from '@driftlock/dataconnect';

// The `ListDetectionTasksByUser` query requires an argument of type `ListDetectionTasksByUserVariables`:
const listDetectionTasksByUserVars: ListDetectionTasksByUserVariables = {
  userId: ..., 
};

// Call the `listDetectionTasksByUserRef()` function to get a reference to the query.
const ref = listDetectionTasksByUserRef(listDetectionTasksByUserVars);
// Variables can be defined inline as well.
const ref = listDetectionTasksByUserRef({ userId: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = listDetectionTasksByUserRef(dataConnect, listDetectionTasksByUserVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.detectionTasks);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.detectionTasks);
});
```

## GetAnomaliesByTask
You can execute the `GetAnomaliesByTask` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
getAnomaliesByTask(vars: GetAnomaliesByTaskVariables): QueryPromise<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;

interface GetAnomaliesByTaskRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetAnomaliesByTaskVariables): QueryRef<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;
}
export const getAnomaliesByTaskRef: GetAnomaliesByTaskRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
getAnomaliesByTask(dc: DataConnect, vars: GetAnomaliesByTaskVariables): QueryPromise<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;

interface GetAnomaliesByTaskRef {
  ...
  (dc: DataConnect, vars: GetAnomaliesByTaskVariables): QueryRef<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;
}
export const getAnomaliesByTaskRef: GetAnomaliesByTaskRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the getAnomaliesByTaskRef:
```typescript
const name = getAnomaliesByTaskRef.operationName;
console.log(name);
```

### Variables
The `GetAnomaliesByTask` query requires an argument of type `GetAnomaliesByTaskVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface GetAnomaliesByTaskVariables {
  taskId: UUIDString;
}
```
### Return Type
Recall that executing the `GetAnomaliesByTask` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `GetAnomaliesByTaskData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface GetAnomaliesByTaskData {
  anomalies: ({
    id: UUIDString;
    dataPointIdentifier: string;
    anomalyScore: number;
    isAnomaly: boolean;
    explanation?: string | null;
    timestamp?: TimestampString | null;
  } & Anomaly_Key)[];
}
```
### Using `GetAnomaliesByTask`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, getAnomaliesByTask, GetAnomaliesByTaskVariables } from '@driftlock/dataconnect';

// The `GetAnomaliesByTask` query requires an argument of type `GetAnomaliesByTaskVariables`:
const getAnomaliesByTaskVars: GetAnomaliesByTaskVariables = {
  taskId: ..., 
};

// Call the `getAnomaliesByTask()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await getAnomaliesByTask(getAnomaliesByTaskVars);
// Variables can be defined inline as well.
const { data } = await getAnomaliesByTask({ taskId: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await getAnomaliesByTask(dataConnect, getAnomaliesByTaskVars);

console.log(data.anomalies);

// Or, you can use the `Promise` API.
getAnomaliesByTask(getAnomaliesByTaskVars).then((response) => {
  const data = response.data;
  console.log(data.anomalies);
});
```

### Using `GetAnomaliesByTask`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, getAnomaliesByTaskRef, GetAnomaliesByTaskVariables } from '@driftlock/dataconnect';

// The `GetAnomaliesByTask` query requires an argument of type `GetAnomaliesByTaskVariables`:
const getAnomaliesByTaskVars: GetAnomaliesByTaskVariables = {
  taskId: ..., 
};

// Call the `getAnomaliesByTaskRef()` function to get a reference to the query.
const ref = getAnomaliesByTaskRef(getAnomaliesByTaskVars);
// Variables can be defined inline as well.
const ref = getAnomaliesByTaskRef({ taskId: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = getAnomaliesByTaskRef(dataConnect, getAnomaliesByTaskVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.anomalies);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.anomalies);
});
```

## GetHighScoreAnomalies
You can execute the `GetHighScoreAnomalies` query using the following action shortcut function, or by calling `executeQuery()` after calling the following `QueryRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
getHighScoreAnomalies(vars: GetHighScoreAnomaliesVariables): QueryPromise<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;

interface GetHighScoreAnomaliesRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetHighScoreAnomaliesVariables): QueryRef<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;
}
export const getHighScoreAnomaliesRef: GetHighScoreAnomaliesRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `QueryRef` function.
```typescript
getHighScoreAnomalies(dc: DataConnect, vars: GetHighScoreAnomaliesVariables): QueryPromise<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;

interface GetHighScoreAnomaliesRef {
  ...
  (dc: DataConnect, vars: GetHighScoreAnomaliesVariables): QueryRef<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;
}
export const getHighScoreAnomaliesRef: GetHighScoreAnomaliesRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the getHighScoreAnomaliesRef:
```typescript
const name = getHighScoreAnomaliesRef.operationName;
console.log(name);
```

### Variables
The `GetHighScoreAnomalies` query requires an argument of type `GetHighScoreAnomaliesVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface GetHighScoreAnomaliesVariables {
  taskId: UUIDString;
  minScore: number;
}
```
### Return Type
Recall that executing the `GetHighScoreAnomalies` query returns a `QueryPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `GetHighScoreAnomaliesData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface GetHighScoreAnomaliesData {
  anomalies: ({
    id: UUIDString;
    dataPointIdentifier: string;
    anomalyScore: number;
    isAnomaly: boolean;
    explanation?: string | null;
    timestamp?: TimestampString | null;
  } & Anomaly_Key)[];
}
```
### Using `GetHighScoreAnomalies`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, getHighScoreAnomalies, GetHighScoreAnomaliesVariables } from '@driftlock/dataconnect';

// The `GetHighScoreAnomalies` query requires an argument of type `GetHighScoreAnomaliesVariables`:
const getHighScoreAnomaliesVars: GetHighScoreAnomaliesVariables = {
  taskId: ..., 
  minScore: ..., 
};

// Call the `getHighScoreAnomalies()` function to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await getHighScoreAnomalies(getHighScoreAnomaliesVars);
// Variables can be defined inline as well.
const { data } = await getHighScoreAnomalies({ taskId: ..., minScore: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await getHighScoreAnomalies(dataConnect, getHighScoreAnomaliesVars);

console.log(data.anomalies);

// Or, you can use the `Promise` API.
getHighScoreAnomalies(getHighScoreAnomaliesVars).then((response) => {
  const data = response.data;
  console.log(data.anomalies);
});
```

### Using `GetHighScoreAnomalies`'s `QueryRef` function

```typescript
import { getDataConnect, executeQuery } from 'firebase/data-connect';
import { connectorConfig, getHighScoreAnomaliesRef, GetHighScoreAnomaliesVariables } from '@driftlock/dataconnect';

// The `GetHighScoreAnomalies` query requires an argument of type `GetHighScoreAnomaliesVariables`:
const getHighScoreAnomaliesVars: GetHighScoreAnomaliesVariables = {
  taskId: ..., 
  minScore: ..., 
};

// Call the `getHighScoreAnomaliesRef()` function to get a reference to the query.
const ref = getHighScoreAnomaliesRef(getHighScoreAnomaliesVars);
// Variables can be defined inline as well.
const ref = getHighScoreAnomaliesRef({ taskId: ..., minScore: ..., });

// You can also pass in a `DataConnect` instance to the `QueryRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = getHighScoreAnomaliesRef(dataConnect, getHighScoreAnomaliesVars);

// Call `executeQuery()` on the reference to execute the query.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeQuery(ref);

console.log(data.anomalies);

// Or, you can use the `Promise` API.
executeQuery(ref).then((response) => {
  const data = response.data;
  console.log(data.anomalies);
});
```

# Mutations

There are two ways to execute a Data Connect Mutation using the generated Web SDK:
- Using a Mutation Reference function, which returns a `MutationRef`
  - The `MutationRef` can be used as an argument to `executeMutation()`, which will execute the Mutation and return a `MutationPromise`
- Using an action shortcut function, which returns a `MutationPromise`
  - Calling the action shortcut function will execute the Mutation and return a `MutationPromise`

The following is true for both the action shortcut function and the `MutationRef` function:
- The `MutationPromise` returned will resolve to the result of the Mutation once it has finished executing
- If the Mutation accepts arguments, both the action shortcut function and the `MutationRef` function accept a single argument: an object that contains all the required variables (and the optional variables) for the Mutation
- Both functions can be called with or without passing in a `DataConnect` instance as an argument. If no `DataConnect` argument is passed in, then the generated SDK will call `getDataConnect(connectorConfig)` behind the scenes for you.

Below are examples of how to use the `driftlock` connector's generated functions to execute each mutation. You can also follow the examples from the [Data Connect documentation](https://firebase.google.com/docs/data-connect/web-sdk#using-mutations).

## CreateUser
You can execute the `CreateUser` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
createUser(vars: CreateUserVariables): MutationPromise<CreateUserData, CreateUserVariables>;

interface CreateUserRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateUserVariables): MutationRef<CreateUserData, CreateUserVariables>;
}
export const createUserRef: CreateUserRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
createUser(dc: DataConnect, vars: CreateUserVariables): MutationPromise<CreateUserData, CreateUserVariables>;

interface CreateUserRef {
  ...
  (dc: DataConnect, vars: CreateUserVariables): MutationRef<CreateUserData, CreateUserVariables>;
}
export const createUserRef: CreateUserRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the createUserRef:
```typescript
const name = createUserRef.operationName;
console.log(name);
```

### Variables
The `CreateUser` mutation requires an argument of type `CreateUserVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface CreateUserVariables {
  displayName: string;
  email?: string | null;
  photoUrl?: string | null;
  createdAt: TimestampString;
}
```
### Return Type
Recall that executing the `CreateUser` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `CreateUserData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface CreateUserData {
  user_insert: User_Key;
}
```
### Using `CreateUser`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, createUser, CreateUserVariables } from '@driftlock/dataconnect';

// The `CreateUser` mutation requires an argument of type `CreateUserVariables`:
const createUserVars: CreateUserVariables = {
  displayName: ..., 
  email: ..., // optional
  photoUrl: ..., // optional
  createdAt: ..., 
};

// Call the `createUser()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await createUser(createUserVars);
// Variables can be defined inline as well.
const { data } = await createUser({ displayName: ..., email: ..., photoUrl: ..., createdAt: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await createUser(dataConnect, createUserVars);

console.log(data.user_insert);

// Or, you can use the `Promise` API.
createUser(createUserVars).then((response) => {
  const data = response.data;
  console.log(data.user_insert);
});
```

### Using `CreateUser`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, createUserRef, CreateUserVariables } from '@driftlock/dataconnect';

// The `CreateUser` mutation requires an argument of type `CreateUserVariables`:
const createUserVars: CreateUserVariables = {
  displayName: ..., 
  email: ..., // optional
  photoUrl: ..., // optional
  createdAt: ..., 
};

// Call the `createUserRef()` function to get a reference to the mutation.
const ref = createUserRef(createUserVars);
// Variables can be defined inline as well.
const ref = createUserRef({ displayName: ..., email: ..., photoUrl: ..., createdAt: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = createUserRef(dataConnect, createUserVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.user_insert);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.user_insert);
});
```

## UpdateUser
You can execute the `UpdateUser` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
updateUser(vars: UpdateUserVariables): MutationPromise<UpdateUserData, UpdateUserVariables>;

interface UpdateUserRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateUserVariables): MutationRef<UpdateUserData, UpdateUserVariables>;
}
export const updateUserRef: UpdateUserRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
updateUser(dc: DataConnect, vars: UpdateUserVariables): MutationPromise<UpdateUserData, UpdateUserVariables>;

interface UpdateUserRef {
  ...
  (dc: DataConnect, vars: UpdateUserVariables): MutationRef<UpdateUserData, UpdateUserVariables>;
}
export const updateUserRef: UpdateUserRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the updateUserRef:
```typescript
const name = updateUserRef.operationName;
console.log(name);
```

### Variables
The `UpdateUser` mutation requires an argument of type `UpdateUserVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface UpdateUserVariables {
  id: UUIDString;
  displayName?: string | null;
  email?: string | null;
  photoUrl?: string | null;
}
```
### Return Type
Recall that executing the `UpdateUser` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `UpdateUserData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface UpdateUserData {
  user_update?: User_Key | null;
}
```
### Using `UpdateUser`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, updateUser, UpdateUserVariables } from '@driftlock/dataconnect';

// The `UpdateUser` mutation requires an argument of type `UpdateUserVariables`:
const updateUserVars: UpdateUserVariables = {
  id: ..., 
  displayName: ..., // optional
  email: ..., // optional
  photoUrl: ..., // optional
};

// Call the `updateUser()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await updateUser(updateUserVars);
// Variables can be defined inline as well.
const { data } = await updateUser({ id: ..., displayName: ..., email: ..., photoUrl: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await updateUser(dataConnect, updateUserVars);

console.log(data.user_update);

// Or, you can use the `Promise` API.
updateUser(updateUserVars).then((response) => {
  const data = response.data;
  console.log(data.user_update);
});
```

### Using `UpdateUser`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, updateUserRef, UpdateUserVariables } from '@driftlock/dataconnect';

// The `UpdateUser` mutation requires an argument of type `UpdateUserVariables`:
const updateUserVars: UpdateUserVariables = {
  id: ..., 
  displayName: ..., // optional
  email: ..., // optional
  photoUrl: ..., // optional
};

// Call the `updateUserRef()` function to get a reference to the mutation.
const ref = updateUserRef(updateUserVars);
// Variables can be defined inline as well.
const ref = updateUserRef({ id: ..., displayName: ..., email: ..., photoUrl: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = updateUserRef(dataConnect, updateUserVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.user_update);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.user_update);
});
```

## CreateDataset
You can execute the `CreateDataset` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
createDataset(vars: CreateDatasetVariables): MutationPromise<CreateDatasetData, CreateDatasetVariables>;

interface CreateDatasetRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateDatasetVariables): MutationRef<CreateDatasetData, CreateDatasetVariables>;
}
export const createDatasetRef: CreateDatasetRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
createDataset(dc: DataConnect, vars: CreateDatasetVariables): MutationPromise<CreateDatasetData, CreateDatasetVariables>;

interface CreateDatasetRef {
  ...
  (dc: DataConnect, vars: CreateDatasetVariables): MutationRef<CreateDatasetData, CreateDatasetVariables>;
}
export const createDatasetRef: CreateDatasetRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the createDatasetRef:
```typescript
const name = createDatasetRef.operationName;
console.log(name);
```

### Variables
The `CreateDataset` mutation requires an argument of type `CreateDatasetVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface CreateDatasetVariables {
  userId: UUIDString;
  name: string;
  filename: string;
  fileLocation: string;
  status: string;
  uploadDate: TimestampString;
  description?: string | null;
  columnInfo?: string | null;
}
```
### Return Type
Recall that executing the `CreateDataset` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `CreateDatasetData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface CreateDatasetData {
  dataset_insert: Dataset_Key;
}
```
### Using `CreateDataset`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, createDataset, CreateDatasetVariables } from '@driftlock/dataconnect';

// The `CreateDataset` mutation requires an argument of type `CreateDatasetVariables`:
const createDatasetVars: CreateDatasetVariables = {
  userId: ..., 
  name: ..., 
  filename: ..., 
  fileLocation: ..., 
  status: ..., 
  uploadDate: ..., 
  description: ..., // optional
  columnInfo: ..., // optional
};

// Call the `createDataset()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await createDataset(createDatasetVars);
// Variables can be defined inline as well.
const { data } = await createDataset({ userId: ..., name: ..., filename: ..., fileLocation: ..., status: ..., uploadDate: ..., description: ..., columnInfo: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await createDataset(dataConnect, createDatasetVars);

console.log(data.dataset_insert);

// Or, you can use the `Promise` API.
createDataset(createDatasetVars).then((response) => {
  const data = response.data;
  console.log(data.dataset_insert);
});
```

### Using `CreateDataset`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, createDatasetRef, CreateDatasetVariables } from '@driftlock/dataconnect';

// The `CreateDataset` mutation requires an argument of type `CreateDatasetVariables`:
const createDatasetVars: CreateDatasetVariables = {
  userId: ..., 
  name: ..., 
  filename: ..., 
  fileLocation: ..., 
  status: ..., 
  uploadDate: ..., 
  description: ..., // optional
  columnInfo: ..., // optional
};

// Call the `createDatasetRef()` function to get a reference to the mutation.
const ref = createDatasetRef(createDatasetVars);
// Variables can be defined inline as well.
const ref = createDatasetRef({ userId: ..., name: ..., filename: ..., fileLocation: ..., status: ..., uploadDate: ..., description: ..., columnInfo: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = createDatasetRef(dataConnect, createDatasetVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.dataset_insert);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.dataset_insert);
});
```

## UpdateDatasetStatus
You can execute the `UpdateDatasetStatus` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
updateDatasetStatus(vars: UpdateDatasetStatusVariables): MutationPromise<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;

interface UpdateDatasetStatusRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateDatasetStatusVariables): MutationRef<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;
}
export const updateDatasetStatusRef: UpdateDatasetStatusRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
updateDatasetStatus(dc: DataConnect, vars: UpdateDatasetStatusVariables): MutationPromise<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;

interface UpdateDatasetStatusRef {
  ...
  (dc: DataConnect, vars: UpdateDatasetStatusVariables): MutationRef<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;
}
export const updateDatasetStatusRef: UpdateDatasetStatusRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the updateDatasetStatusRef:
```typescript
const name = updateDatasetStatusRef.operationName;
console.log(name);
```

### Variables
The `UpdateDatasetStatus` mutation requires an argument of type `UpdateDatasetStatusVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface UpdateDatasetStatusVariables {
  id: UUIDString;
  status: string;
}
```
### Return Type
Recall that executing the `UpdateDatasetStatus` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `UpdateDatasetStatusData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface UpdateDatasetStatusData {
  dataset_update?: Dataset_Key | null;
}
```
### Using `UpdateDatasetStatus`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, updateDatasetStatus, UpdateDatasetStatusVariables } from '@driftlock/dataconnect';

// The `UpdateDatasetStatus` mutation requires an argument of type `UpdateDatasetStatusVariables`:
const updateDatasetStatusVars: UpdateDatasetStatusVariables = {
  id: ..., 
  status: ..., 
};

// Call the `updateDatasetStatus()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await updateDatasetStatus(updateDatasetStatusVars);
// Variables can be defined inline as well.
const { data } = await updateDatasetStatus({ id: ..., status: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await updateDatasetStatus(dataConnect, updateDatasetStatusVars);

console.log(data.dataset_update);

// Or, you can use the `Promise` API.
updateDatasetStatus(updateDatasetStatusVars).then((response) => {
  const data = response.data;
  console.log(data.dataset_update);
});
```

### Using `UpdateDatasetStatus`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, updateDatasetStatusRef, UpdateDatasetStatusVariables } from '@driftlock/dataconnect';

// The `UpdateDatasetStatus` mutation requires an argument of type `UpdateDatasetStatusVariables`:
const updateDatasetStatusVars: UpdateDatasetStatusVariables = {
  id: ..., 
  status: ..., 
};

// Call the `updateDatasetStatusRef()` function to get a reference to the mutation.
const ref = updateDatasetStatusRef(updateDatasetStatusVars);
// Variables can be defined inline as well.
const ref = updateDatasetStatusRef({ id: ..., status: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = updateDatasetStatusRef(dataConnect, updateDatasetStatusVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.dataset_update);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.dataset_update);
});
```

## UpdateDataset
You can execute the `UpdateDataset` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
updateDataset(vars: UpdateDatasetVariables): MutationPromise<UpdateDatasetData, UpdateDatasetVariables>;

interface UpdateDatasetRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateDatasetVariables): MutationRef<UpdateDatasetData, UpdateDatasetVariables>;
}
export const updateDatasetRef: UpdateDatasetRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
updateDataset(dc: DataConnect, vars: UpdateDatasetVariables): MutationPromise<UpdateDatasetData, UpdateDatasetVariables>;

interface UpdateDatasetRef {
  ...
  (dc: DataConnect, vars: UpdateDatasetVariables): MutationRef<UpdateDatasetData, UpdateDatasetVariables>;
}
export const updateDatasetRef: UpdateDatasetRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the updateDatasetRef:
```typescript
const name = updateDatasetRef.operationName;
console.log(name);
```

### Variables
The `UpdateDataset` mutation requires an argument of type `UpdateDatasetVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface UpdateDatasetVariables {
  id: UUIDString;
  name?: string | null;
  description?: string | null;
  status?: string | null;
  columnInfo?: string | null;
}
```
### Return Type
Recall that executing the `UpdateDataset` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `UpdateDatasetData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface UpdateDatasetData {
  dataset_update?: Dataset_Key | null;
}
```
### Using `UpdateDataset`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, updateDataset, UpdateDatasetVariables } from '@driftlock/dataconnect';

// The `UpdateDataset` mutation requires an argument of type `UpdateDatasetVariables`:
const updateDatasetVars: UpdateDatasetVariables = {
  id: ..., 
  name: ..., // optional
  description: ..., // optional
  status: ..., // optional
  columnInfo: ..., // optional
};

// Call the `updateDataset()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await updateDataset(updateDatasetVars);
// Variables can be defined inline as well.
const { data } = await updateDataset({ id: ..., name: ..., description: ..., status: ..., columnInfo: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await updateDataset(dataConnect, updateDatasetVars);

console.log(data.dataset_update);

// Or, you can use the `Promise` API.
updateDataset(updateDatasetVars).then((response) => {
  const data = response.data;
  console.log(data.dataset_update);
});
```

### Using `UpdateDataset`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, updateDatasetRef, UpdateDatasetVariables } from '@driftlock/dataconnect';

// The `UpdateDataset` mutation requires an argument of type `UpdateDatasetVariables`:
const updateDatasetVars: UpdateDatasetVariables = {
  id: ..., 
  name: ..., // optional
  description: ..., // optional
  status: ..., // optional
  columnInfo: ..., // optional
};

// Call the `updateDatasetRef()` function to get a reference to the mutation.
const ref = updateDatasetRef(updateDatasetVars);
// Variables can be defined inline as well.
const ref = updateDatasetRef({ id: ..., name: ..., description: ..., status: ..., columnInfo: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = updateDatasetRef(dataConnect, updateDatasetVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.dataset_update);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.dataset_update);
});
```

## CreateModelConfiguration
You can execute the `CreateModelConfiguration` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
createModelConfiguration(vars: CreateModelConfigurationVariables): MutationPromise<CreateModelConfigurationData, CreateModelConfigurationVariables>;

interface CreateModelConfigurationRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateModelConfigurationVariables): MutationRef<CreateModelConfigurationData, CreateModelConfigurationVariables>;
}
export const createModelConfigurationRef: CreateModelConfigurationRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
createModelConfiguration(dc: DataConnect, vars: CreateModelConfigurationVariables): MutationPromise<CreateModelConfigurationData, CreateModelConfigurationVariables>;

interface CreateModelConfigurationRef {
  ...
  (dc: DataConnect, vars: CreateModelConfigurationVariables): MutationRef<CreateModelConfigurationData, CreateModelConfigurationVariables>;
}
export const createModelConfigurationRef: CreateModelConfigurationRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the createModelConfigurationRef:
```typescript
const name = createModelConfigurationRef.operationName;
console.log(name);
```

### Variables
The `CreateModelConfiguration` mutation requires an argument of type `CreateModelConfigurationVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface CreateModelConfigurationVariables {
  userId: UUIDString;
  name: string;
  algorithmType: string;
  createdAt: TimestampString;
  parameters?: string | null;
  description?: string | null;
}
```
### Return Type
Recall that executing the `CreateModelConfiguration` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `CreateModelConfigurationData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface CreateModelConfigurationData {
  modelConfiguration_insert: ModelConfiguration_Key;
}
```
### Using `CreateModelConfiguration`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, createModelConfiguration, CreateModelConfigurationVariables } from '@driftlock/dataconnect';

// The `CreateModelConfiguration` mutation requires an argument of type `CreateModelConfigurationVariables`:
const createModelConfigurationVars: CreateModelConfigurationVariables = {
  userId: ..., 
  name: ..., 
  algorithmType: ..., 
  createdAt: ..., 
  parameters: ..., // optional
  description: ..., // optional
};

// Call the `createModelConfiguration()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await createModelConfiguration(createModelConfigurationVars);
// Variables can be defined inline as well.
const { data } = await createModelConfiguration({ userId: ..., name: ..., algorithmType: ..., createdAt: ..., parameters: ..., description: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await createModelConfiguration(dataConnect, createModelConfigurationVars);

console.log(data.modelConfiguration_insert);

// Or, you can use the `Promise` API.
createModelConfiguration(createModelConfigurationVars).then((response) => {
  const data = response.data;
  console.log(data.modelConfiguration_insert);
});
```

### Using `CreateModelConfiguration`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, createModelConfigurationRef, CreateModelConfigurationVariables } from '@driftlock/dataconnect';

// The `CreateModelConfiguration` mutation requires an argument of type `CreateModelConfigurationVariables`:
const createModelConfigurationVars: CreateModelConfigurationVariables = {
  userId: ..., 
  name: ..., 
  algorithmType: ..., 
  createdAt: ..., 
  parameters: ..., // optional
  description: ..., // optional
};

// Call the `createModelConfigurationRef()` function to get a reference to the mutation.
const ref = createModelConfigurationRef(createModelConfigurationVars);
// Variables can be defined inline as well.
const ref = createModelConfigurationRef({ userId: ..., name: ..., algorithmType: ..., createdAt: ..., parameters: ..., description: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = createModelConfigurationRef(dataConnect, createModelConfigurationVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.modelConfiguration_insert);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.modelConfiguration_insert);
});
```

## UpdateModelConfiguration
You can execute the `UpdateModelConfiguration` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
updateModelConfiguration(vars: UpdateModelConfigurationVariables): MutationPromise<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;

interface UpdateModelConfigurationRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateModelConfigurationVariables): MutationRef<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;
}
export const updateModelConfigurationRef: UpdateModelConfigurationRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
updateModelConfiguration(dc: DataConnect, vars: UpdateModelConfigurationVariables): MutationPromise<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;

interface UpdateModelConfigurationRef {
  ...
  (dc: DataConnect, vars: UpdateModelConfigurationVariables): MutationRef<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;
}
export const updateModelConfigurationRef: UpdateModelConfigurationRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the updateModelConfigurationRef:
```typescript
const name = updateModelConfigurationRef.operationName;
console.log(name);
```

### Variables
The `UpdateModelConfiguration` mutation requires an argument of type `UpdateModelConfigurationVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface UpdateModelConfigurationVariables {
  id: UUIDString;
  name?: string | null;
  parameters?: string | null;
  description?: string | null;
}
```
### Return Type
Recall that executing the `UpdateModelConfiguration` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `UpdateModelConfigurationData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface UpdateModelConfigurationData {
  modelConfiguration_update?: ModelConfiguration_Key | null;
}
```
### Using `UpdateModelConfiguration`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, updateModelConfiguration, UpdateModelConfigurationVariables } from '@driftlock/dataconnect';

// The `UpdateModelConfiguration` mutation requires an argument of type `UpdateModelConfigurationVariables`:
const updateModelConfigurationVars: UpdateModelConfigurationVariables = {
  id: ..., 
  name: ..., // optional
  parameters: ..., // optional
  description: ..., // optional
};

// Call the `updateModelConfiguration()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await updateModelConfiguration(updateModelConfigurationVars);
// Variables can be defined inline as well.
const { data } = await updateModelConfiguration({ id: ..., name: ..., parameters: ..., description: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await updateModelConfiguration(dataConnect, updateModelConfigurationVars);

console.log(data.modelConfiguration_update);

// Or, you can use the `Promise` API.
updateModelConfiguration(updateModelConfigurationVars).then((response) => {
  const data = response.data;
  console.log(data.modelConfiguration_update);
});
```

### Using `UpdateModelConfiguration`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, updateModelConfigurationRef, UpdateModelConfigurationVariables } from '@driftlock/dataconnect';

// The `UpdateModelConfiguration` mutation requires an argument of type `UpdateModelConfigurationVariables`:
const updateModelConfigurationVars: UpdateModelConfigurationVariables = {
  id: ..., 
  name: ..., // optional
  parameters: ..., // optional
  description: ..., // optional
};

// Call the `updateModelConfigurationRef()` function to get a reference to the mutation.
const ref = updateModelConfigurationRef(updateModelConfigurationVars);
// Variables can be defined inline as well.
const ref = updateModelConfigurationRef({ id: ..., name: ..., parameters: ..., description: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = updateModelConfigurationRef(dataConnect, updateModelConfigurationVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.modelConfiguration_update);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.modelConfiguration_update);
});
```

## CreateDetectionTask
You can execute the `CreateDetectionTask` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
createDetectionTask(vars: CreateDetectionTaskVariables): MutationPromise<CreateDetectionTaskData, CreateDetectionTaskVariables>;

interface CreateDetectionTaskRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateDetectionTaskVariables): MutationRef<CreateDetectionTaskData, CreateDetectionTaskVariables>;
}
export const createDetectionTaskRef: CreateDetectionTaskRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
createDetectionTask(dc: DataConnect, vars: CreateDetectionTaskVariables): MutationPromise<CreateDetectionTaskData, CreateDetectionTaskVariables>;

interface CreateDetectionTaskRef {
  ...
  (dc: DataConnect, vars: CreateDetectionTaskVariables): MutationRef<CreateDetectionTaskData, CreateDetectionTaskVariables>;
}
export const createDetectionTaskRef: CreateDetectionTaskRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the createDetectionTaskRef:
```typescript
const name = createDetectionTaskRef.operationName;
console.log(name);
```

### Variables
The `CreateDetectionTask` mutation requires an argument of type `CreateDetectionTaskVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface CreateDetectionTaskVariables {
  userId: UUIDString;
  datasetId: UUIDString;
  modelConfigurationId: UUIDString;
  taskName: string;
  status: string;
  startDate: TimestampString;
  notes?: string | null;
}
```
### Return Type
Recall that executing the `CreateDetectionTask` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `CreateDetectionTaskData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface CreateDetectionTaskData {
  detectionTask_insert: DetectionTask_Key;
}
```
### Using `CreateDetectionTask`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, createDetectionTask, CreateDetectionTaskVariables } from '@driftlock/dataconnect';

// The `CreateDetectionTask` mutation requires an argument of type `CreateDetectionTaskVariables`:
const createDetectionTaskVars: CreateDetectionTaskVariables = {
  userId: ..., 
  datasetId: ..., 
  modelConfigurationId: ..., 
  taskName: ..., 
  status: ..., 
  startDate: ..., 
  notes: ..., // optional
};

// Call the `createDetectionTask()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await createDetectionTask(createDetectionTaskVars);
// Variables can be defined inline as well.
const { data } = await createDetectionTask({ userId: ..., datasetId: ..., modelConfigurationId: ..., taskName: ..., status: ..., startDate: ..., notes: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await createDetectionTask(dataConnect, createDetectionTaskVars);

console.log(data.detectionTask_insert);

// Or, you can use the `Promise` API.
createDetectionTask(createDetectionTaskVars).then((response) => {
  const data = response.data;
  console.log(data.detectionTask_insert);
});
```

### Using `CreateDetectionTask`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, createDetectionTaskRef, CreateDetectionTaskVariables } from '@driftlock/dataconnect';

// The `CreateDetectionTask` mutation requires an argument of type `CreateDetectionTaskVariables`:
const createDetectionTaskVars: CreateDetectionTaskVariables = {
  userId: ..., 
  datasetId: ..., 
  modelConfigurationId: ..., 
  taskName: ..., 
  status: ..., 
  startDate: ..., 
  notes: ..., // optional
};

// Call the `createDetectionTaskRef()` function to get a reference to the mutation.
const ref = createDetectionTaskRef(createDetectionTaskVars);
// Variables can be defined inline as well.
const ref = createDetectionTaskRef({ userId: ..., datasetId: ..., modelConfigurationId: ..., taskName: ..., status: ..., startDate: ..., notes: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = createDetectionTaskRef(dataConnect, createDetectionTaskVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.detectionTask_insert);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.detectionTask_insert);
});
```

## UpdateDetectionTask
You can execute the `UpdateDetectionTask` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
updateDetectionTask(vars: UpdateDetectionTaskVariables): MutationPromise<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;

interface UpdateDetectionTaskRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateDetectionTaskVariables): MutationRef<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;
}
export const updateDetectionTaskRef: UpdateDetectionTaskRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
updateDetectionTask(dc: DataConnect, vars: UpdateDetectionTaskVariables): MutationPromise<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;

interface UpdateDetectionTaskRef {
  ...
  (dc: DataConnect, vars: UpdateDetectionTaskVariables): MutationRef<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;
}
export const updateDetectionTaskRef: UpdateDetectionTaskRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the updateDetectionTaskRef:
```typescript
const name = updateDetectionTaskRef.operationName;
console.log(name);
```

### Variables
The `UpdateDetectionTask` mutation requires an argument of type `UpdateDetectionTaskVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface UpdateDetectionTaskVariables {
  id: UUIDString;
  status?: string | null;
  resultsLocation?: string | null;
  notes?: string | null;
  endDate?: TimestampString | null;
}
```
### Return Type
Recall that executing the `UpdateDetectionTask` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `UpdateDetectionTaskData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface UpdateDetectionTaskData {
  detectionTask_update?: DetectionTask_Key | null;
}
```
### Using `UpdateDetectionTask`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, updateDetectionTask, UpdateDetectionTaskVariables } from '@driftlock/dataconnect';

// The `UpdateDetectionTask` mutation requires an argument of type `UpdateDetectionTaskVariables`:
const updateDetectionTaskVars: UpdateDetectionTaskVariables = {
  id: ..., 
  status: ..., // optional
  resultsLocation: ..., // optional
  notes: ..., // optional
  endDate: ..., // optional
};

// Call the `updateDetectionTask()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await updateDetectionTask(updateDetectionTaskVars);
// Variables can be defined inline as well.
const { data } = await updateDetectionTask({ id: ..., status: ..., resultsLocation: ..., notes: ..., endDate: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await updateDetectionTask(dataConnect, updateDetectionTaskVars);

console.log(data.detectionTask_update);

// Or, you can use the `Promise` API.
updateDetectionTask(updateDetectionTaskVars).then((response) => {
  const data = response.data;
  console.log(data.detectionTask_update);
});
```

### Using `UpdateDetectionTask`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, updateDetectionTaskRef, UpdateDetectionTaskVariables } from '@driftlock/dataconnect';

// The `UpdateDetectionTask` mutation requires an argument of type `UpdateDetectionTaskVariables`:
const updateDetectionTaskVars: UpdateDetectionTaskVariables = {
  id: ..., 
  status: ..., // optional
  resultsLocation: ..., // optional
  notes: ..., // optional
  endDate: ..., // optional
};

// Call the `updateDetectionTaskRef()` function to get a reference to the mutation.
const ref = updateDetectionTaskRef(updateDetectionTaskVars);
// Variables can be defined inline as well.
const ref = updateDetectionTaskRef({ id: ..., status: ..., resultsLocation: ..., notes: ..., endDate: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = updateDetectionTaskRef(dataConnect, updateDetectionTaskVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.detectionTask_update);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.detectionTask_update);
});
```

## CompleteDetectionTask
You can execute the `CompleteDetectionTask` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
completeDetectionTask(vars: CompleteDetectionTaskVariables): MutationPromise<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;

interface CompleteDetectionTaskRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: CompleteDetectionTaskVariables): MutationRef<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;
}
export const completeDetectionTaskRef: CompleteDetectionTaskRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
completeDetectionTask(dc: DataConnect, vars: CompleteDetectionTaskVariables): MutationPromise<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;

interface CompleteDetectionTaskRef {
  ...
  (dc: DataConnect, vars: CompleteDetectionTaskVariables): MutationRef<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;
}
export const completeDetectionTaskRef: CompleteDetectionTaskRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the completeDetectionTaskRef:
```typescript
const name = completeDetectionTaskRef.operationName;
console.log(name);
```

### Variables
The `CompleteDetectionTask` mutation requires an argument of type `CompleteDetectionTaskVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface CompleteDetectionTaskVariables {
  id: UUIDString;
  resultsLocation: string;
  endDate: TimestampString;
}
```
### Return Type
Recall that executing the `CompleteDetectionTask` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `CompleteDetectionTaskData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface CompleteDetectionTaskData {
  detectionTask_update?: DetectionTask_Key | null;
}
```
### Using `CompleteDetectionTask`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, completeDetectionTask, CompleteDetectionTaskVariables } from '@driftlock/dataconnect';

// The `CompleteDetectionTask` mutation requires an argument of type `CompleteDetectionTaskVariables`:
const completeDetectionTaskVars: CompleteDetectionTaskVariables = {
  id: ..., 
  resultsLocation: ..., 
  endDate: ..., 
};

// Call the `completeDetectionTask()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await completeDetectionTask(completeDetectionTaskVars);
// Variables can be defined inline as well.
const { data } = await completeDetectionTask({ id: ..., resultsLocation: ..., endDate: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await completeDetectionTask(dataConnect, completeDetectionTaskVars);

console.log(data.detectionTask_update);

// Or, you can use the `Promise` API.
completeDetectionTask(completeDetectionTaskVars).then((response) => {
  const data = response.data;
  console.log(data.detectionTask_update);
});
```

### Using `CompleteDetectionTask`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, completeDetectionTaskRef, CompleteDetectionTaskVariables } from '@driftlock/dataconnect';

// The `CompleteDetectionTask` mutation requires an argument of type `CompleteDetectionTaskVariables`:
const completeDetectionTaskVars: CompleteDetectionTaskVariables = {
  id: ..., 
  resultsLocation: ..., 
  endDate: ..., 
};

// Call the `completeDetectionTaskRef()` function to get a reference to the mutation.
const ref = completeDetectionTaskRef(completeDetectionTaskVars);
// Variables can be defined inline as well.
const ref = completeDetectionTaskRef({ id: ..., resultsLocation: ..., endDate: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = completeDetectionTaskRef(dataConnect, completeDetectionTaskVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.detectionTask_update);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.detectionTask_update);
});
```

## CreateAnomaly
You can execute the `CreateAnomaly` mutation using the following action shortcut function, or by calling `executeMutation()` after calling the following `MutationRef` function, both of which are defined in [generated/index.d.ts](./index.d.ts):
```typescript
createAnomaly(vars: CreateAnomalyVariables): MutationPromise<CreateAnomalyData, CreateAnomalyVariables>;

interface CreateAnomalyRef {
  ...
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateAnomalyVariables): MutationRef<CreateAnomalyData, CreateAnomalyVariables>;
}
export const createAnomalyRef: CreateAnomalyRef;
```
You can also pass in a `DataConnect` instance to the action shortcut function or `MutationRef` function.
```typescript
createAnomaly(dc: DataConnect, vars: CreateAnomalyVariables): MutationPromise<CreateAnomalyData, CreateAnomalyVariables>;

interface CreateAnomalyRef {
  ...
  (dc: DataConnect, vars: CreateAnomalyVariables): MutationRef<CreateAnomalyData, CreateAnomalyVariables>;
}
export const createAnomalyRef: CreateAnomalyRef;
```

If you need the name of the operation without creating a ref, you can retrieve the operation name by calling the `operationName` property on the createAnomalyRef:
```typescript
const name = createAnomalyRef.operationName;
console.log(name);
```

### Variables
The `CreateAnomaly` mutation requires an argument of type `CreateAnomalyVariables`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:

```typescript
export interface CreateAnomalyVariables {
  detectionTaskId: UUIDString;
  dataPointIdentifier: string;
  anomalyScore: number;
  isAnomaly: boolean;
  explanation?: string | null;
  timestamp?: TimestampString | null;
}
```
### Return Type
Recall that executing the `CreateAnomaly` mutation returns a `MutationPromise` that resolves to an object with a `data` property.

The `data` property is an object of type `CreateAnomalyData`, which is defined in [generated/index.d.ts](./index.d.ts). It has the following fields:
```typescript
export interface CreateAnomalyData {
  anomaly_insert: Anomaly_Key;
}
```
### Using `CreateAnomaly`'s action shortcut function

```typescript
import { getDataConnect } from 'firebase/data-connect';
import { connectorConfig, createAnomaly, CreateAnomalyVariables } from '@driftlock/dataconnect';

// The `CreateAnomaly` mutation requires an argument of type `CreateAnomalyVariables`:
const createAnomalyVars: CreateAnomalyVariables = {
  detectionTaskId: ..., 
  dataPointIdentifier: ..., 
  anomalyScore: ..., 
  isAnomaly: ..., 
  explanation: ..., // optional
  timestamp: ..., // optional
};

// Call the `createAnomaly()` function to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await createAnomaly(createAnomalyVars);
// Variables can be defined inline as well.
const { data } = await createAnomaly({ detectionTaskId: ..., dataPointIdentifier: ..., anomalyScore: ..., isAnomaly: ..., explanation: ..., timestamp: ..., });

// You can also pass in a `DataConnect` instance to the action shortcut function.
const dataConnect = getDataConnect(connectorConfig);
const { data } = await createAnomaly(dataConnect, createAnomalyVars);

console.log(data.anomaly_insert);

// Or, you can use the `Promise` API.
createAnomaly(createAnomalyVars).then((response) => {
  const data = response.data;
  console.log(data.anomaly_insert);
});
```

### Using `CreateAnomaly`'s `MutationRef` function

```typescript
import { getDataConnect, executeMutation } from 'firebase/data-connect';
import { connectorConfig, createAnomalyRef, CreateAnomalyVariables } from '@driftlock/dataconnect';

// The `CreateAnomaly` mutation requires an argument of type `CreateAnomalyVariables`:
const createAnomalyVars: CreateAnomalyVariables = {
  detectionTaskId: ..., 
  dataPointIdentifier: ..., 
  anomalyScore: ..., 
  isAnomaly: ..., 
  explanation: ..., // optional
  timestamp: ..., // optional
};

// Call the `createAnomalyRef()` function to get a reference to the mutation.
const ref = createAnomalyRef(createAnomalyVars);
// Variables can be defined inline as well.
const ref = createAnomalyRef({ detectionTaskId: ..., dataPointIdentifier: ..., anomalyScore: ..., isAnomaly: ..., explanation: ..., timestamp: ..., });

// You can also pass in a `DataConnect` instance to the `MutationRef` function.
const dataConnect = getDataConnect(connectorConfig);
const ref = createAnomalyRef(dataConnect, createAnomalyVars);

// Call `executeMutation()` on the reference to execute the mutation.
// You can use the `await` keyword to wait for the promise to resolve.
const { data } = await executeMutation(ref);

console.log(data.anomaly_insert);

// Or, you can use the `Promise` API.
executeMutation(ref).then((response) => {
  const data = response.data;
  console.log(data.anomaly_insert);
});
```

