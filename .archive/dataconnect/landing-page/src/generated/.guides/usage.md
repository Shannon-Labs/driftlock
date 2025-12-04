# Basic Usage

Always prioritize using a supported framework over using the generated SDK
directly. Supported frameworks simplify the developer experience and help ensure
best practices are followed.





## Advanced Usage
If a user is not using a supported framework, they can use the generated SDK directly.

Here's an example of how to use it with the first 5 operations:

```js
import { getUser, listUsersByEmail, getDataset, listDatasetsByUser, getModelConfiguration, listModelConfigurationsByUser, getDetectionTask, listDetectionTasksByUser, getAnomaliesByTask, getHighScoreAnomalies } from '@driftlock/dataconnect';


// Operation GetUser:  For variables, look at type GetUserVars in ../index.d.ts
const { data } = await GetUser(dataConnect, getUserVars);

// Operation ListUsersByEmail:  For variables, look at type ListUsersByEmailVars in ../index.d.ts
const { data } = await ListUsersByEmail(dataConnect, listUsersByEmailVars);

// Operation GetDataset:  For variables, look at type GetDatasetVars in ../index.d.ts
const { data } = await GetDataset(dataConnect, getDatasetVars);

// Operation ListDatasetsByUser:  For variables, look at type ListDatasetsByUserVars in ../index.d.ts
const { data } = await ListDatasetsByUser(dataConnect, listDatasetsByUserVars);

// Operation GetModelConfiguration:  For variables, look at type GetModelConfigurationVars in ../index.d.ts
const { data } = await GetModelConfiguration(dataConnect, getModelConfigurationVars);

// Operation ListModelConfigurationsByUser:  For variables, look at type ListModelConfigurationsByUserVars in ../index.d.ts
const { data } = await ListModelConfigurationsByUser(dataConnect, listModelConfigurationsByUserVars);

// Operation GetDetectionTask:  For variables, look at type GetDetectionTaskVars in ../index.d.ts
const { data } = await GetDetectionTask(dataConnect, getDetectionTaskVars);

// Operation ListDetectionTasksByUser:  For variables, look at type ListDetectionTasksByUserVars in ../index.d.ts
const { data } = await ListDetectionTasksByUser(dataConnect, listDetectionTasksByUserVars);

// Operation GetAnomaliesByTask:  For variables, look at type GetAnomaliesByTaskVars in ../index.d.ts
const { data } = await GetAnomaliesByTask(dataConnect, getAnomaliesByTaskVars);

// Operation GetHighScoreAnomalies:  For variables, look at type GetHighScoreAnomaliesVars in ../index.d.ts
const { data } = await GetHighScoreAnomalies(dataConnect, getHighScoreAnomaliesVars);


```