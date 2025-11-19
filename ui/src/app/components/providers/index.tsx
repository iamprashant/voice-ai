import { Metadata } from '@rapidaai/react';

export interface ProviderConfig {
  provider: string;
  parameters: Metadata[];
}

export interface ProviderComponentProps {
  provider: string;
  onChangeProvider: (provider: string) => void;
  parameters: Metadata[];
  onChangeParameter: (parameters: Metadata[]) => void;
}
