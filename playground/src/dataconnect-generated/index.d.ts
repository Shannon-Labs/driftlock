import { ConnectorConfig, DataConnect, QueryRef, QueryPromise, MutationRef, MutationPromise } from 'firebase/data-connect';

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

interface CreateUserRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateUserVariables): MutationRef<CreateUserData, CreateUserVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: CreateUserVariables): MutationRef<CreateUserData, CreateUserVariables>;
  operationName: string;
}
export const createUserRef: CreateUserRef;

export function createUser(vars: CreateUserVariables): MutationPromise<CreateUserData, CreateUserVariables>;
export function createUser(dc: DataConnect, vars: CreateUserVariables): MutationPromise<CreateUserData, CreateUserVariables>;

interface UpdateUserRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateUserVariables): MutationRef<UpdateUserData, UpdateUserVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: UpdateUserVariables): MutationRef<UpdateUserData, UpdateUserVariables>;
  operationName: string;
}
export const updateUserRef: UpdateUserRef;

export function updateUser(vars: UpdateUserVariables): MutationPromise<UpdateUserData, UpdateUserVariables>;
export function updateUser(dc: DataConnect, vars: UpdateUserVariables): MutationPromise<UpdateUserData, UpdateUserVariables>;

interface CreateDatasetRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateDatasetVariables): MutationRef<CreateDatasetData, CreateDatasetVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: CreateDatasetVariables): MutationRef<CreateDatasetData, CreateDatasetVariables>;
  operationName: string;
}
export const createDatasetRef: CreateDatasetRef;

export function createDataset(vars: CreateDatasetVariables): MutationPromise<CreateDatasetData, CreateDatasetVariables>;
export function createDataset(dc: DataConnect, vars: CreateDatasetVariables): MutationPromise<CreateDatasetData, CreateDatasetVariables>;

interface UpdateDatasetStatusRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateDatasetStatusVariables): MutationRef<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: UpdateDatasetStatusVariables): MutationRef<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;
  operationName: string;
}
export const updateDatasetStatusRef: UpdateDatasetStatusRef;

export function updateDatasetStatus(vars: UpdateDatasetStatusVariables): MutationPromise<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;
export function updateDatasetStatus(dc: DataConnect, vars: UpdateDatasetStatusVariables): MutationPromise<UpdateDatasetStatusData, UpdateDatasetStatusVariables>;

interface UpdateDatasetRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateDatasetVariables): MutationRef<UpdateDatasetData, UpdateDatasetVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: UpdateDatasetVariables): MutationRef<UpdateDatasetData, UpdateDatasetVariables>;
  operationName: string;
}
export const updateDatasetRef: UpdateDatasetRef;

export function updateDataset(vars: UpdateDatasetVariables): MutationPromise<UpdateDatasetData, UpdateDatasetVariables>;
export function updateDataset(dc: DataConnect, vars: UpdateDatasetVariables): MutationPromise<UpdateDatasetData, UpdateDatasetVariables>;

interface CreateModelConfigurationRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateModelConfigurationVariables): MutationRef<CreateModelConfigurationData, CreateModelConfigurationVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: CreateModelConfigurationVariables): MutationRef<CreateModelConfigurationData, CreateModelConfigurationVariables>;
  operationName: string;
}
export const createModelConfigurationRef: CreateModelConfigurationRef;

export function createModelConfiguration(vars: CreateModelConfigurationVariables): MutationPromise<CreateModelConfigurationData, CreateModelConfigurationVariables>;
export function createModelConfiguration(dc: DataConnect, vars: CreateModelConfigurationVariables): MutationPromise<CreateModelConfigurationData, CreateModelConfigurationVariables>;

interface UpdateModelConfigurationRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateModelConfigurationVariables): MutationRef<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: UpdateModelConfigurationVariables): MutationRef<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;
  operationName: string;
}
export const updateModelConfigurationRef: UpdateModelConfigurationRef;

export function updateModelConfiguration(vars: UpdateModelConfigurationVariables): MutationPromise<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;
export function updateModelConfiguration(dc: DataConnect, vars: UpdateModelConfigurationVariables): MutationPromise<UpdateModelConfigurationData, UpdateModelConfigurationVariables>;

interface CreateDetectionTaskRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateDetectionTaskVariables): MutationRef<CreateDetectionTaskData, CreateDetectionTaskVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: CreateDetectionTaskVariables): MutationRef<CreateDetectionTaskData, CreateDetectionTaskVariables>;
  operationName: string;
}
export const createDetectionTaskRef: CreateDetectionTaskRef;

export function createDetectionTask(vars: CreateDetectionTaskVariables): MutationPromise<CreateDetectionTaskData, CreateDetectionTaskVariables>;
export function createDetectionTask(dc: DataConnect, vars: CreateDetectionTaskVariables): MutationPromise<CreateDetectionTaskData, CreateDetectionTaskVariables>;

interface UpdateDetectionTaskRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: UpdateDetectionTaskVariables): MutationRef<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: UpdateDetectionTaskVariables): MutationRef<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;
  operationName: string;
}
export const updateDetectionTaskRef: UpdateDetectionTaskRef;

export function updateDetectionTask(vars: UpdateDetectionTaskVariables): MutationPromise<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;
export function updateDetectionTask(dc: DataConnect, vars: UpdateDetectionTaskVariables): MutationPromise<UpdateDetectionTaskData, UpdateDetectionTaskVariables>;

interface CompleteDetectionTaskRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: CompleteDetectionTaskVariables): MutationRef<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: CompleteDetectionTaskVariables): MutationRef<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;
  operationName: string;
}
export const completeDetectionTaskRef: CompleteDetectionTaskRef;

export function completeDetectionTask(vars: CompleteDetectionTaskVariables): MutationPromise<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;
export function completeDetectionTask(dc: DataConnect, vars: CompleteDetectionTaskVariables): MutationPromise<CompleteDetectionTaskData, CompleteDetectionTaskVariables>;

interface CreateAnomalyRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: CreateAnomalyVariables): MutationRef<CreateAnomalyData, CreateAnomalyVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: CreateAnomalyVariables): MutationRef<CreateAnomalyData, CreateAnomalyVariables>;
  operationName: string;
}
export const createAnomalyRef: CreateAnomalyRef;

export function createAnomaly(vars: CreateAnomalyVariables): MutationPromise<CreateAnomalyData, CreateAnomalyVariables>;
export function createAnomaly(dc: DataConnect, vars: CreateAnomalyVariables): MutationPromise<CreateAnomalyData, CreateAnomalyVariables>;

interface GetUserRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetUserVariables): QueryRef<GetUserData, GetUserVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: GetUserVariables): QueryRef<GetUserData, GetUserVariables>;
  operationName: string;
}
export const getUserRef: GetUserRef;

export function getUser(vars: GetUserVariables): QueryPromise<GetUserData, GetUserVariables>;
export function getUser(dc: DataConnect, vars: GetUserVariables): QueryPromise<GetUserData, GetUserVariables>;

interface ListUsersByEmailRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListUsersByEmailVariables): QueryRef<ListUsersByEmailData, ListUsersByEmailVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: ListUsersByEmailVariables): QueryRef<ListUsersByEmailData, ListUsersByEmailVariables>;
  operationName: string;
}
export const listUsersByEmailRef: ListUsersByEmailRef;

export function listUsersByEmail(vars: ListUsersByEmailVariables): QueryPromise<ListUsersByEmailData, ListUsersByEmailVariables>;
export function listUsersByEmail(dc: DataConnect, vars: ListUsersByEmailVariables): QueryPromise<ListUsersByEmailData, ListUsersByEmailVariables>;

interface GetDatasetRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetDatasetVariables): QueryRef<GetDatasetData, GetDatasetVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: GetDatasetVariables): QueryRef<GetDatasetData, GetDatasetVariables>;
  operationName: string;
}
export const getDatasetRef: GetDatasetRef;

export function getDataset(vars: GetDatasetVariables): QueryPromise<GetDatasetData, GetDatasetVariables>;
export function getDataset(dc: DataConnect, vars: GetDatasetVariables): QueryPromise<GetDatasetData, GetDatasetVariables>;

interface ListDatasetsByUserRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListDatasetsByUserVariables): QueryRef<ListDatasetsByUserData, ListDatasetsByUserVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: ListDatasetsByUserVariables): QueryRef<ListDatasetsByUserData, ListDatasetsByUserVariables>;
  operationName: string;
}
export const listDatasetsByUserRef: ListDatasetsByUserRef;

export function listDatasetsByUser(vars: ListDatasetsByUserVariables): QueryPromise<ListDatasetsByUserData, ListDatasetsByUserVariables>;
export function listDatasetsByUser(dc: DataConnect, vars: ListDatasetsByUserVariables): QueryPromise<ListDatasetsByUserData, ListDatasetsByUserVariables>;

interface GetModelConfigurationRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetModelConfigurationVariables): QueryRef<GetModelConfigurationData, GetModelConfigurationVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: GetModelConfigurationVariables): QueryRef<GetModelConfigurationData, GetModelConfigurationVariables>;
  operationName: string;
}
export const getModelConfigurationRef: GetModelConfigurationRef;

export function getModelConfiguration(vars: GetModelConfigurationVariables): QueryPromise<GetModelConfigurationData, GetModelConfigurationVariables>;
export function getModelConfiguration(dc: DataConnect, vars: GetModelConfigurationVariables): QueryPromise<GetModelConfigurationData, GetModelConfigurationVariables>;

interface ListModelConfigurationsByUserRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListModelConfigurationsByUserVariables): QueryRef<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: ListModelConfigurationsByUserVariables): QueryRef<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;
  operationName: string;
}
export const listModelConfigurationsByUserRef: ListModelConfigurationsByUserRef;

export function listModelConfigurationsByUser(vars: ListModelConfigurationsByUserVariables): QueryPromise<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;
export function listModelConfigurationsByUser(dc: DataConnect, vars: ListModelConfigurationsByUserVariables): QueryPromise<ListModelConfigurationsByUserData, ListModelConfigurationsByUserVariables>;

interface GetDetectionTaskRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetDetectionTaskVariables): QueryRef<GetDetectionTaskData, GetDetectionTaskVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: GetDetectionTaskVariables): QueryRef<GetDetectionTaskData, GetDetectionTaskVariables>;
  operationName: string;
}
export const getDetectionTaskRef: GetDetectionTaskRef;

export function getDetectionTask(vars: GetDetectionTaskVariables): QueryPromise<GetDetectionTaskData, GetDetectionTaskVariables>;
export function getDetectionTask(dc: DataConnect, vars: GetDetectionTaskVariables): QueryPromise<GetDetectionTaskData, GetDetectionTaskVariables>;

interface ListDetectionTasksByUserRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: ListDetectionTasksByUserVariables): QueryRef<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: ListDetectionTasksByUserVariables): QueryRef<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;
  operationName: string;
}
export const listDetectionTasksByUserRef: ListDetectionTasksByUserRef;

export function listDetectionTasksByUser(vars: ListDetectionTasksByUserVariables): QueryPromise<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;
export function listDetectionTasksByUser(dc: DataConnect, vars: ListDetectionTasksByUserVariables): QueryPromise<ListDetectionTasksByUserData, ListDetectionTasksByUserVariables>;

interface GetAnomaliesByTaskRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetAnomaliesByTaskVariables): QueryRef<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: GetAnomaliesByTaskVariables): QueryRef<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;
  operationName: string;
}
export const getAnomaliesByTaskRef: GetAnomaliesByTaskRef;

export function getAnomaliesByTask(vars: GetAnomaliesByTaskVariables): QueryPromise<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;
export function getAnomaliesByTask(dc: DataConnect, vars: GetAnomaliesByTaskVariables): QueryPromise<GetAnomaliesByTaskData, GetAnomaliesByTaskVariables>;

interface GetHighScoreAnomaliesRef {
  /* Allow users to create refs without passing in DataConnect */
  (vars: GetHighScoreAnomaliesVariables): QueryRef<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;
  /* Allow users to pass in custom DataConnect instances */
  (dc: DataConnect, vars: GetHighScoreAnomaliesVariables): QueryRef<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;
  operationName: string;
}
export const getHighScoreAnomaliesRef: GetHighScoreAnomaliesRef;

export function getHighScoreAnomalies(vars: GetHighScoreAnomaliesVariables): QueryPromise<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;
export function getHighScoreAnomalies(dc: DataConnect, vars: GetHighScoreAnomaliesVariables): QueryPromise<GetHighScoreAnomaliesData, GetHighScoreAnomaliesVariables>;

