import { Metadata } from '@rapidaai/react';

export interface ProviderComponentProps {
  provider: string;
  onChangeProvider: (provider: string) => void;
  parameters: Metadata[];
  onChangeParameter: (parameters: Metadata[]) => void;
}
