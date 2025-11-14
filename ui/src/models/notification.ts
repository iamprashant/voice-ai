export const RAPIDA_SYSTEM_NOTIFICATION = [
  {
    category: 'Assistant',
    items: [
      {
        id: 'assistant.created',
        label: 'Assistant Created',
        description: 'Triggered when a new assistant is created.',
        default: true,
      },
      {
        id: 'assistant.updated',
        label: 'Assistant Updated',
        description:
          'Triggered when assistant configuration or metadata are updated.',
        default: true,
      },
      {
        id: 'assistant.version.deployed',
        label: 'New Version Deployed',
        description:
          'Notifies when a new version of the assistant is deployed.',
        default: true,
      },
      {
        id: 'assistant.version.rollback',
        label: 'Version Rolled Back',
        description:
          'Triggered when an assistant is rolled back to a previous version.',
        default: true,
      },
      {
        id: 'assistant.deleted',
        label: 'Assistant Deleted',
        description: 'Triggered when an assistant is deleted.',
        default: true,
      },
      {
        id: 'assistant.deployment.created',
        label: 'Deployment Created',
        description:
          'Triggered when a new deployment is created for an assistant.',
        default: true,
      },
      {
        id: 'assistant.deployment.updated',
        label: 'Deployment Updated',
        description:
          'Triggered when an assistant deployment configuration is updated.',
        default: true,
      },
      {
        id: 'assistant.deployment.failed',
        label: 'Deployment Failed',
        description:
          'Triggered when a deployment fails to complete successfully.',
        default: true,
      },
      {
        id: 'assistant.deployment.succeeded',
        label: 'Deployment Succeeded',
        description:
          'Triggered when an assistant deployment completes successfully.',
        default: true,
      },
      {
        id: 'assistant.error',
        label: 'Assistant Error Detected',
        description:
          'Triggered when an assistant experiences repeated errors or degraded performance.',
        default: true,
      },
    ],
  },
  {
    category: 'Knowledge',
    items: [
      {
        id: 'knowledge.source.added',
        label: 'Knowledge Source Added',
        description:
          'Triggered when a new knowledge source is added (document, website, or dataset).',
        default: true,
      },
      {
        id: 'knowledge.ingest.started',
        label: 'Knowledge Ingestion Started',
        description: 'Triggered when ingestion of a knowledge source begins.',
        default: false,
      },
      {
        id: 'knowledge.ingest.completed',
        label: 'Knowledge Ingestion Completed',
        description:
          'Triggered when a knowledge ingestion job completes successfully.',
        default: true,
      },
      {
        id: 'knowledge.ingest.failed',
        label: 'Knowledge Ingestion Failed',
        description: 'Triggered when a knowledge ingestion job fails.',
        default: true,
      },
      {
        id: 'knowledge.source.deleted',
        label: 'Knowledge Source Deleted',
        description:
          'Triggered when a knowledge source is removed from the system.',
        default: true,
      },
    ],
  },
  {
    category: 'Endpoint',
    items: [
      {
        id: 'endpoint.created',
        label: 'Endpoint Created',
        description: 'Triggered when a new LLM or API endpoint is created.',
        default: true,
      },
      {
        id: 'endpoint.updated',
        label: 'Endpoint Updated',
        description:
          'Triggered when endpoint configuration or model settings are updated.',
        default: true,
      },
      {
        id: 'endpoint.failed',
        label: 'Endpoint Failure',
        description:
          'Triggered when endpoint requests repeatedly fail or time out.',
        default: true,
      },
      {
        id: 'endpoint.performance.degraded',
        label: 'Endpoint Performance Degraded',
        description:
          'Triggered when latency or response errors exceed threshold.',
        default: true,
      },
      {
        id: 'endpoint.deleted',
        label: 'Endpoint Deleted',
        description: 'Triggered when an endpoint is deleted or disabled.',
        default: true,
      },
    ],
  },
  {
    category: 'Integration',
    items: [
      {
        id: 'integration.connected',
        label: 'Integration Connected',
        description:
          'Triggered when a third-party integration (e.g., Slack, SendGrid) is successfully connected.',
        default: true,
      },
      {
        id: 'integration.disconnected',
        label: 'Integration Disconnected',
        description:
          'Triggered when an integration is disconnected or authorization expires.',
        default: true,
      },
      {
        id: 'integration.failed',
        label: 'Integration Failure',
        description:
          'Triggered when an integration call fails (auth error, webhook timeout, etc.).',
        default: true,
      },
      {
        id: 'integration.updated',
        label: 'Integration Updated',
        description:
          'Triggered when integration settings or credentials are modified.',
        default: true,
      },
    ],
  },
];
