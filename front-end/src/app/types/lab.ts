

export interface ContainerStatusInfo {
    name: string;
    ready: boolean;
    state: string;
    reason: string;
    startedAt: string;
  }

  export interface GetIDEStatusResponse {
    status: string;
    lastHeartbeat: string;
  }


  export interface CreateLabResponse {
    ideUrl: string;
    namespace: string;
    expiresAt: string;
    labPassword: string;
  }

  export interface FinishLabResponse {
    message: string;
  }

  export interface GetAllLabsByUserAndStatusResponse {
    namespace: string;
    labSlug: string;
    startedAt: string;
    finishedAt: string;
    score: number;
    status: string;
    durationSeconds: number;
    rating: number;
    userFeedback: string;
  }

  export interface FinishResult {
    status: string;
    errorMessage: string;
    totalScore: number;
    criteriaResult: LabFinishResultCriterion[];
    filesDiff: string;
  }

  export interface LabFinishResultCriterion {
    name: string;
    score: number;
    weight: number;
    required: boolean;
    message: string;
    status: string;
    rawOutput: string;
  }

  export interface GetLabResultResponse {
    namespace: string;
    labSlug: string;
    status: string;
    rating: number;
    userFeedback: string;
    startedAt: string;
    finishedAt: string;
    durationSeconds: number;
    finishResult: FinishResult;
  }
  
  
  export interface LabDefinition {
    id: string;
    slug: string;
    title: string;
    authors: string[];
    externalReferences: string[];
    estimatedTime: number;
    difficulty: 'easy' | 'medium' | 'hard';
    description: string;
    readme: string;
    tags: string[];
  }


export interface LabVulnerability {
  name: string;
  description: string;
  owaspReference: string;
}

export interface LabLanguage {
  name: string;
}

export interface LabTechnology {
  name: string;
}

export interface Lab {
  id: string;
  title: string;
  description: string;
  vulnerabilities: LabVulnerability[];
  difficulty: string;
  languages: LabLanguage[];
  technologies: LabTechnology[];
  references: string[];
  authors: string[];
  rating: number;
  estimatedTime: number;
  slug: string;
  createdAt: string;
  updatedAt: string;
}