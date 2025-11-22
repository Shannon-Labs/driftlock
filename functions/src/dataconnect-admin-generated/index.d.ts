import { ConnectorConfig, DataConnect, OperationOptions, ExecuteOperationResponse } from 'firebase-admin/data-connect';

export const connectorConfig: ConnectorConfig;

export type TimestampString = string;
export type UUIDString = string;
export type Int64String = string;
export type DateString = string;


export interface Anomaly_Key {
  id: UUIDString;
  __typename?: 'Anomaly_Key';
}

export interface CompleteDetectionTaskData {
  detectionTask_update?: DetectionTask_Key | null;
}

export interface CompleteDetectionTaskVariables {
  id: UUIDString;
  resultsLocation: string;
  endDate: TimestampString;
}

export interface CreateAnomalyData {
  anomaly_insert: Anomaly_Key;
}

export interface CreateAnomalyVariables {
  detectionTaskId: UUIDString;
  dataPointIdentifier: string;
  anomalyScore: number;
  isAnomaly: boolean;
  explanation?: string | null;
  timestamp?: TimestampString | null;
}

export interface CreateDatasetData {
  dataset_insert: Dataset_Key;
}

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

export interface CreateDetectionTaskData {
  detectionTask_insert: DetectionTask_Key;
}

export interface CreateDetectionTaskVariables {
  userId: UUIDString;
  datasetId: UUIDString;
  modelConfigurationId: UUIDString;
  taskName: string;
  status: string;
  startDate: TimestampString;
  notes?: string | null;
}

export interface CreateModelConfigurationData {
  modelConfiguration_insert: ModelConfiguration_Key;
}

export interface CreateModelConfigurationVariables {
  userId: UUIDString;
  name: string;
  algorithmType: string;
  createdAt: TimestampString;
  parameters?: string | null;
  description?: string | null;
}

export interface CreateUserData {
  user_insert: User_Key;
}

export interface CreateUserVariables {
  displayName: string;
  email?: string | null;
  photoUrl?: string | null;
  createdAt: TimestampString;
}

export interface Dataset_Key {
  id: UUIDString;
  __typename?: 'Dataset_Key';
}

export interface DetectionTask_Key {
  id: UUIDString;
  __typename?: 'DetectionTask_Key';
}

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

export interface GetAnomaliesByTaskVariables {
  taskId: UUIDString;
}

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

export interface GetDatasetVariables {
  id: UUIDString;
}

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

export interface GetDetectionTaskVariables {
  id: UUIDString;
}

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

export interface GetHighScoreAnomaliesVariables {
  taskId: UUIDString;
  minScore: number;
}

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

export interface GetModelConfigurationVariables {
  id: UUIDString;
}

export interface GetUserData {
  user?: {
    id: UUIDString;
    displayName: string;
    email?: string | null;
    photoUrl?: string | null;
    createdAt: TimestampString;
  } & User_Key;
}

export interface GetUserVariables {
  id: UUIDString;
}

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

export interface ListDatasetsByUserVariables {
  userId: UUIDString;
}

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

export interface ListDetectionTasksByUserVariables {
  userId: UUIDString;
}

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

export interface ListModelConfigurationsByUserVariables {
  userId: UUIDString;
}

export interface ListUsersByEmailData {
  users: ({
    id: UUIDString;
    displayName: string;
    email?: string | null;
    photoUrl?: string | null;
    createdAt: TimestampString;
  } & User_Key)[];
}

export interface ListUsersByEmailVariables {
  email: string;
}

export interface ModelConfiguration_Key {
  id: UUIDString;
  __typename?: 'ModelConfiguration_Key';
}

export interface UpdateDatasetData {
  dataset_update?: Dataset_Key | null;
}

export interface UpdateDatasetStatusData {
  dataset_update?: Dataset_Key | null;
}

export interface UpdateDatasetStatusVariables {
  id: UUIDString;
  status: string;
}

export interface UpdateDatasetVariables {
  id: UUIDString;
  name?: string | null;
  description?: string | null;
  status?: string | null;
  columnInfo?: string | null;
}

export interface UpdateDetectionTaskData {
  detectionTask_update?: DetectionTask_Key | null;
}

export interface UpdateDetectionTaskVariables {
  id: UUIDString;
  status?: string | null;
  resultsLocation?: string | null;
  notes?: string | null;
  endDate?: TimestampString | null;
}

export interface UpdateModelConfigurationData {
  modelConfiguration_update?: ModelConfiguration_Key | null;
}

export interface UpdateModelConfigurationVariables {
  id: UUIDString;
  name?: string | null;
  parameters?: string | null;
  description?: string | null;
}

export interface UpdateUserData {
  user_update?: User_Key | null;
}

export interface UpdateUserVariables {
  id: UUIDString;
  displayName?: string | null;
  email?: string | null;
  photoUrl?: string | null;
}

export interface User_Key {
  id: UUIDString;
  __typename?: 'User_Key';
}

/** Generated Node Admin SDK operation action function for the 'CreateUser' Mutation. Allow users to execute without passing in DataConnect. */
export function createUser(dc: DataConnect, vars: CreateUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateUserData>>;
/** Generated Node Admin SDK operation action function for the 'CreateUser' Mutation. Allow users to pass in custom DataConnect instances. */
export function createUser(vars: CreateUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateUserData>>;

/** Generated Node Admin SDK operation action function for the 'UpdateUser' Mutation. Allow users to execute without passing in DataConnect. */
export function updateUser(dc: DataConnect, vars: UpdateUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateUserData>>;
/** Generated Node Admin SDK operation action function for the 'UpdateUser' Mutation. Allow users to pass in custom DataConnect instances. */
export function updateUser(vars: UpdateUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateUserData>>;

/** Generated Node Admin SDK operation action function for the 'CreateDataset' Mutation. Allow users to execute without passing in DataConnect. */
export function createDataset(dc: DataConnect, vars: CreateDatasetVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateDatasetData>>;
/** Generated Node Admin SDK operation action function for the 'CreateDataset' Mutation. Allow users to pass in custom DataConnect instances. */
export function createDataset(vars: CreateDatasetVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateDatasetData>>;

/** Generated Node Admin SDK operation action function for the 'UpdateDatasetStatus' Mutation. Allow users to execute without passing in DataConnect. */
export function updateDatasetStatus(dc: DataConnect, vars: UpdateDatasetStatusVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateDatasetStatusData>>;
/** Generated Node Admin SDK operation action function for the 'UpdateDatasetStatus' Mutation. Allow users to pass in custom DataConnect instances. */
export function updateDatasetStatus(vars: UpdateDatasetStatusVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateDatasetStatusData>>;

/** Generated Node Admin SDK operation action function for the 'UpdateDataset' Mutation. Allow users to execute without passing in DataConnect. */
export function updateDataset(dc: DataConnect, vars: UpdateDatasetVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateDatasetData>>;
/** Generated Node Admin SDK operation action function for the 'UpdateDataset' Mutation. Allow users to pass in custom DataConnect instances. */
export function updateDataset(vars: UpdateDatasetVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateDatasetData>>;

/** Generated Node Admin SDK operation action function for the 'CreateModelConfiguration' Mutation. Allow users to execute without passing in DataConnect. */
export function createModelConfiguration(dc: DataConnect, vars: CreateModelConfigurationVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateModelConfigurationData>>;
/** Generated Node Admin SDK operation action function for the 'CreateModelConfiguration' Mutation. Allow users to pass in custom DataConnect instances. */
export function createModelConfiguration(vars: CreateModelConfigurationVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateModelConfigurationData>>;

/** Generated Node Admin SDK operation action function for the 'UpdateModelConfiguration' Mutation. Allow users to execute without passing in DataConnect. */
export function updateModelConfiguration(dc: DataConnect, vars: UpdateModelConfigurationVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateModelConfigurationData>>;
/** Generated Node Admin SDK operation action function for the 'UpdateModelConfiguration' Mutation. Allow users to pass in custom DataConnect instances. */
export function updateModelConfiguration(vars: UpdateModelConfigurationVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateModelConfigurationData>>;

/** Generated Node Admin SDK operation action function for the 'CreateDetectionTask' Mutation. Allow users to execute without passing in DataConnect. */
export function createDetectionTask(dc: DataConnect, vars: CreateDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateDetectionTaskData>>;
/** Generated Node Admin SDK operation action function for the 'CreateDetectionTask' Mutation. Allow users to pass in custom DataConnect instances. */
export function createDetectionTask(vars: CreateDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateDetectionTaskData>>;

/** Generated Node Admin SDK operation action function for the 'UpdateDetectionTask' Mutation. Allow users to execute without passing in DataConnect. */
export function updateDetectionTask(dc: DataConnect, vars: UpdateDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateDetectionTaskData>>;
/** Generated Node Admin SDK operation action function for the 'UpdateDetectionTask' Mutation. Allow users to pass in custom DataConnect instances. */
export function updateDetectionTask(vars: UpdateDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<UpdateDetectionTaskData>>;

/** Generated Node Admin SDK operation action function for the 'CompleteDetectionTask' Mutation. Allow users to execute without passing in DataConnect. */
export function completeDetectionTask(dc: DataConnect, vars: CompleteDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CompleteDetectionTaskData>>;
/** Generated Node Admin SDK operation action function for the 'CompleteDetectionTask' Mutation. Allow users to pass in custom DataConnect instances. */
export function completeDetectionTask(vars: CompleteDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CompleteDetectionTaskData>>;

/** Generated Node Admin SDK operation action function for the 'CreateAnomaly' Mutation. Allow users to execute without passing in DataConnect. */
export function createAnomaly(dc: DataConnect, vars: CreateAnomalyVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateAnomalyData>>;
/** Generated Node Admin SDK operation action function for the 'CreateAnomaly' Mutation. Allow users to pass in custom DataConnect instances. */
export function createAnomaly(vars: CreateAnomalyVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<CreateAnomalyData>>;

/** Generated Node Admin SDK operation action function for the 'GetUser' Query. Allow users to execute without passing in DataConnect. */
export function getUser(dc: DataConnect, vars: GetUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetUserData>>;
/** Generated Node Admin SDK operation action function for the 'GetUser' Query. Allow users to pass in custom DataConnect instances. */
export function getUser(vars: GetUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetUserData>>;

/** Generated Node Admin SDK operation action function for the 'ListUsersByEmail' Query. Allow users to execute without passing in DataConnect. */
export function listUsersByEmail(dc: DataConnect, vars: ListUsersByEmailVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListUsersByEmailData>>;
/** Generated Node Admin SDK operation action function for the 'ListUsersByEmail' Query. Allow users to pass in custom DataConnect instances. */
export function listUsersByEmail(vars: ListUsersByEmailVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListUsersByEmailData>>;

/** Generated Node Admin SDK operation action function for the 'GetDataset' Query. Allow users to execute without passing in DataConnect. */
export function getDataset(dc: DataConnect, vars: GetDatasetVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetDatasetData>>;
/** Generated Node Admin SDK operation action function for the 'GetDataset' Query. Allow users to pass in custom DataConnect instances. */
export function getDataset(vars: GetDatasetVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetDatasetData>>;

/** Generated Node Admin SDK operation action function for the 'ListDatasetsByUser' Query. Allow users to execute without passing in DataConnect. */
export function listDatasetsByUser(dc: DataConnect, vars: ListDatasetsByUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListDatasetsByUserData>>;
/** Generated Node Admin SDK operation action function for the 'ListDatasetsByUser' Query. Allow users to pass in custom DataConnect instances. */
export function listDatasetsByUser(vars: ListDatasetsByUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListDatasetsByUserData>>;

/** Generated Node Admin SDK operation action function for the 'GetModelConfiguration' Query. Allow users to execute without passing in DataConnect. */
export function getModelConfiguration(dc: DataConnect, vars: GetModelConfigurationVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetModelConfigurationData>>;
/** Generated Node Admin SDK operation action function for the 'GetModelConfiguration' Query. Allow users to pass in custom DataConnect instances. */
export function getModelConfiguration(vars: GetModelConfigurationVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetModelConfigurationData>>;

/** Generated Node Admin SDK operation action function for the 'ListModelConfigurationsByUser' Query. Allow users to execute without passing in DataConnect. */
export function listModelConfigurationsByUser(dc: DataConnect, vars: ListModelConfigurationsByUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListModelConfigurationsByUserData>>;
/** Generated Node Admin SDK operation action function for the 'ListModelConfigurationsByUser' Query. Allow users to pass in custom DataConnect instances. */
export function listModelConfigurationsByUser(vars: ListModelConfigurationsByUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListModelConfigurationsByUserData>>;

/** Generated Node Admin SDK operation action function for the 'GetDetectionTask' Query. Allow users to execute without passing in DataConnect. */
export function getDetectionTask(dc: DataConnect, vars: GetDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetDetectionTaskData>>;
/** Generated Node Admin SDK operation action function for the 'GetDetectionTask' Query. Allow users to pass in custom DataConnect instances. */
export function getDetectionTask(vars: GetDetectionTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetDetectionTaskData>>;

/** Generated Node Admin SDK operation action function for the 'ListDetectionTasksByUser' Query. Allow users to execute without passing in DataConnect. */
export function listDetectionTasksByUser(dc: DataConnect, vars: ListDetectionTasksByUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListDetectionTasksByUserData>>;
/** Generated Node Admin SDK operation action function for the 'ListDetectionTasksByUser' Query. Allow users to pass in custom DataConnect instances. */
export function listDetectionTasksByUser(vars: ListDetectionTasksByUserVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<ListDetectionTasksByUserData>>;

/** Generated Node Admin SDK operation action function for the 'GetAnomaliesByTask' Query. Allow users to execute without passing in DataConnect. */
export function getAnomaliesByTask(dc: DataConnect, vars: GetAnomaliesByTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetAnomaliesByTaskData>>;
/** Generated Node Admin SDK operation action function for the 'GetAnomaliesByTask' Query. Allow users to pass in custom DataConnect instances. */
export function getAnomaliesByTask(vars: GetAnomaliesByTaskVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetAnomaliesByTaskData>>;

/** Generated Node Admin SDK operation action function for the 'GetHighScoreAnomalies' Query. Allow users to execute without passing in DataConnect. */
export function getHighScoreAnomalies(dc: DataConnect, vars: GetHighScoreAnomaliesVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetHighScoreAnomaliesData>>;
/** Generated Node Admin SDK operation action function for the 'GetHighScoreAnomalies' Query. Allow users to pass in custom DataConnect instances. */
export function getHighScoreAnomalies(vars: GetHighScoreAnomaliesVariables, options?: OperationOptions): Promise<ExecuteOperationResponse<GetHighScoreAnomaliesData>>;

