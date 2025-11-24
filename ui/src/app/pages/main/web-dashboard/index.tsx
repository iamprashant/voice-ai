import {
  ChevronRight,
  ExternalLink,
  MessageSquare,
  TestTube,
} from 'lucide-react';
import { ClickableCard } from '@/app/components/base/cards';
import { useCurrentCredential } from '@/hooks/use-credential';
import { KnowledgeIcon } from '@/app/components/Icon/knowledge';
import { EndpointIcon } from '@/app/components/Icon/Endpoint';
import { AssistantIcon } from '@/app/components/Icon/Assistant';
import { ModelIcon } from '@/app/components/Icon/Model';
import { useProviderContext } from '@/context/provider-context';
import { ILinkButton } from '@/app/components/form/button';

export const HomePage = () => {
  const coreFeatures = [
    {
      icon: KnowledgeIcon,
      title: 'Knowledge Hub',
      description:
        'Unified repository for documents, training data, and AI knowledge management — the foundation of contextual intelligence.',
      color: 'bg-blue-500',
      route: '/knowledge',
    },
    {
      icon: MessageSquare,
      title: 'Conversational AI',
      description:
        'Context-aware, LLM-powered chat experiences that understand user intent and deliver accurate responses.',
      color: 'bg-yellow-500',
      route: '/deployment/assistant',
    },
    {
      icon: AssistantIcon,
      title: 'AI Assistants',
      description:
        'Deploy domain-specific AI agents with custom skills, workflows, and multi-step reasoning.',
      color: 'bg-green-500',
      route: '/deployment/assistant',
    },
    {
      icon: EndpointIcon,
      title: 'Governance & Endpoints',
      description:
        'Secure API endpoints with fine-grained governance, audit trails, and enterprise-grade access control.',
      color: 'bg-purple-500',
      route: '/deployment',
    },
    {
      icon: ModelIcon,
      title: 'Model Integration',
      description:
        'Bring your own model — support for OpenAI, Anthropic, and custom LLMs with fine-tuning capabilities.',
      color: 'bg-red-500',
      route: '/integration/models',
    },
    {
      icon: TestTube,
      title: 'Real-time Testing & Monitoring',
      description:
        'Instantly test AI agents and flows in a live sandbox to iterate faster and ship confidently.',
      color: 'bg-indigo-500',
      route: '/logs',
    },
  ];
  const { user } = useCurrentCredential();
  const { providerCredentials } = useProviderContext();
  return (
    <div className="flex-1 overflow-auto flex flex-col">
      {/* Core Platform Features */}
      {providerCredentials.length === 0 && (
        <>
          <div className="p-4">
            <h2 className="text-xl font-semibold text-foreground mb-1">
              Get Started
            </h2>
            <p className="text-sm text-muted-foreground">
              Complete these steps to set up your workspace
            </p>
          </div>
          <div className="border-y p-3">
            <div className="flex justify-between border bg-light-background dark:bg-gray-950 items-center">
              <div className="p-4">
                <h1 className="text-base font-semibold">
                  Connect your ai providers
                </h1>
                <p className="text-sm text-muted-foreground mt-1">
                  Set up the language model, speech-to-text, and text-to-speech
                  services you want to use.
                </p>
              </div>
              <ILinkButton
                href="/integration/models"
                className="p-0 px-2 pl-4 mr-4"
              >
                Connect
                <ExternalLink className="w-4 h-4 ml-2" strokeWidth={1.5} />
              </ILinkButton>
            </div>
          </div>
        </>
      )}

      <div className="border-b bg-white dark:bg-gray-900 p-4">
        <h1 className="text-xl font-semibold">
          Welcome, {user?.name?.split(/\s+/)[0] || user?.name}{' '}
        </h1>
      </div>
      <main className="px-6 py-6 bg-white dark:bg-gray-900">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
          {coreFeatures.map((feature, index) => (
            <ClickableCard
              to={feature.route}
              key={index}
              className="transition-all duration-200 cursor-pointer group p-4 hover:shadow-lg border-b-2 hover:border-blue-600 h-full border"
            >
              <div className="flex flex-col space-y-4">
                <div className="flex items-center justify-between">
                  <div
                    className={`w-10 h-10 ${feature.color} flex items-center justify-center group-hover:scale-110 transition-transform`}
                  >
                    <feature.icon
                      className="h-5 w-5 text-white"
                      strokeWidth={1.5}
                    />
                  </div>
                </div>
                <div>
                  <h3 className="text-lg font-semibold">{feature.title}</h3>
                  <p className="text-[0.95rem] text-gray-600 dark:text-gray-500 mt-2">
                    {feature.description}
                  </p>
                </div>
              </div>
            </ClickableCard>
          ))}
        </div>
      </main>
      <div className="border-y p-4 justify-between items-center bg-white dark:bg-gray-900">
        <h1 className="">
          Have questions, need support, or want to explore how we can help you
          build better AI experiences? Our team is always ready to assist.
        </h1>
        <p>
          Reach out anytime — Get quick help from our team at:
          <a
            href="mailto:tech@rapida.ai"
            className="mx-2 text-blue-600 hover:underline underline-offset-2"
          >
            tech@rapida.ai
          </a>
        </p>
      </div>
    </div>
  );
};
