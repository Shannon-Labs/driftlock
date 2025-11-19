# Basic Usage

Always prioritize using a supported framework over using the generated SDK
directly. Supported frameworks simplify the developer experience and help ensure
best practices are followed.





## Advanced Usage
If a user is not using a supported framework, they can use the generated SDK directly.

Here's an example of how to use it with the first 5 operations:

```js
import { createUser, updateUser, createDataset, updateDatasetStatus, updateDataset, createModelConfiguration, updateModelConfiguration, createDetectionTask, updateDetectionTask, completeDetectionTask } from '@driftlock/dataconnect';


// Operation CreateUser:  For variables, look at type CreateUserVars in ../index.d.ts
const { data } = await CreateUser(dataConnect, createUserVars);

// Operation UpdateUser:  For variables, look at type UpdateUserVars in ../index.d.ts
const { data } = await UpdateUser(dataConnect, updateUserVars);

// Operation CreateDataset:  For variables, look at type CreateDatasetVars in ../index.d.ts
const { data } = await CreateDataset(dataConnect, createDatasetVars);

// Operation UpdateDatasetStatus:  For variables, look at type UpdateDatasetStatusVars in ../index.d.ts
const { data } = await UpdateDatasetStatus(dataConnect, updateDatasetStatusVars);

// Operation UpdateDataset:  For variables, look at type UpdateDatasetVars in ../index.d.ts
const { data } = await UpdateDataset(dataConnect, updateDatasetVars);

// Operation CreateModelConfiguration:  For variables, look at type CreateModelConfigurationVars in ../index.d.ts
const { data } = await CreateModelConfiguration(dataConnect, createModelConfigurationVars);

// Operation UpdateModelConfiguration:  For variables, look at type UpdateModelConfigurationVars in ../index.d.ts
const { data } = await UpdateModelConfiguration(dataConnect, updateModelConfigurationVars);

// Operation CreateDetectionTask:  For variables, look at type CreateDetectionTaskVars in ../index.d.ts
const { data } = await CreateDetectionTask(dataConnect, createDetectionTaskVars);

// Operation UpdateDetectionTask:  For variables, look at type UpdateDetectionTaskVars in ../index.d.ts
const { data } = await UpdateDetectionTask(dataConnect, updateDetectionTaskVars);

// Operation CompleteDetectionTask:  For variables, look at type CompleteDetectionTaskVars in ../index.d.ts
const { data } = await CompleteDetectionTask(dataConnect, completeDetectionTaskVars);


```