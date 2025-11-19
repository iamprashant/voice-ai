import productionProvider from './provider.production.json';
import developmentProvider from './provider.development.json';
import elevan_lab_voices from './elevanlabs/voices.json';
import azure_voices from './azure/voices.json';
import cartesia_voices from './cartesia/voices.json';
import google_cloud_voices from './google/voices.json';
import deepgram_voices from './deepgram/voices.json';

export interface IntegrationProvider extends RapidaProvider {}
interface EndOfSpeechProvider extends RapidaProvider {}
interface NoiseCancellationProvider extends RapidaProvider {}
export interface RapidaProvider {
  code: string;
  name: string;
  featureList: string[];
  description?: string;
  image?: string;
  url?: string;
  configurations?: {
    name: string;
    type: string;
    label: string;
  }[];
}

export const allProvider = (): RapidaProvider[] => {
  const env = process.env.NODE_ENV || 'development';
  return env === 'production' ? productionProvider : developmentProvider;
};

export const EndOfSpeech = (): EndOfSpeechProvider[] => {
  return allProvider().filter(x => x.featureList.includes('end_of_speech'));
};

export const NoiseCancellation = (): NoiseCancellationProvider[] => {
  return allProvider().filter(x =>
    x.featureList.includes('noise_cancellation'),
  );
};

export const TEXT_PROVIDERS = allProvider().filter(x =>
  x.featureList.includes('text'),
);
export const INTEGRATION_PROVIDER: IntegrationProvider[] = allProvider().filter(
  x => x.featureList.includes('external'),
);
export const SPEECH_TO_TEXT_PROVIDER = allProvider().filter(x =>
  x.featureList.includes('stt'),
);
export const TEXT_TO_SPEECH_PROVIDER = allProvider().filter(x =>
  x.featureList.includes('tts'),
);
export const TEXT_TO_SPEECH = (name: string) =>
  allProvider()
    .filter(x => x.featureList.includes('tts'))
    .findLast(x => x.code === name);

export const STORAGE_PROVIDER = allProvider().filter(x =>
  x.featureList.includes('storage'),
);

export const RERANKER_PROVIDER = allProvider().filter(x =>
  x.featureList.includes('reranker'),
);

export const TELEPHONY_PROVIDER = allProvider().filter(x =>
  x.featureList.includes('telephony'),
);

export const EMBEDDING_PROVIDERS = allProvider().filter(x =>
  x.featureList.includes('embedding'),
);

export const ELEVANLABS_VOICE = () => {
  return elevan_lab_voices;
};

export const AZURE_VOICE = () => {
  return azure_voices;
};

export const CARTESIA_VOICE = () => {
  return cartesia_voices;
};

export const GOOGLE_CLOUD_VOICE = () => {
  return google_cloud_voices;
};

export const DEEPGRAM_VOICE = () => {
  return deepgram_voices;
};
