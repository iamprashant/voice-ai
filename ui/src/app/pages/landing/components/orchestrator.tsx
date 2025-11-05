import { motion } from 'framer-motion';
import { FC, useEffect, useState } from 'react';
import {
  AudioLines,
  BookOpen,
  Bot,
  Cloud,
  Code,
  Globe2,
  Mail,
  MessageCircle,
  MessageSquare,
  Phone,
  PhoneForwarded,
  RefreshCcw,
  TabletSmartphone,
  Workflow,
} from 'lucide-react';
import { cn } from '@/styles/media';

type AnimationSteps = {
  animateChannel: boolean;
  animateSpeechToText: boolean;
  animateLLM: boolean;
  animateTool: boolean;
  animateTextToSpeech: boolean;
};
export const Orchestrator: React.FC = () => {
  const [steps, setSteps] = useState<AnimationSteps>({
    animateChannel: false,
    animateSpeechToText: false,
    animateLLM: false,
    animateTool: false,
    animateTextToSpeech: false,
  });

  useEffect(() => {
    const sequence = async () => {
      // Step 1: Channel
      setSteps({
        animateChannel: true,
        animateSpeechToText: false,
        animateLLM: false,
        animateTool: false,
        animateTextToSpeech: false,
      });
      await new Promise(res => setTimeout(res, 3000));

      // Step 2: SpeechToText
      setSteps({
        animateChannel: false,
        animateSpeechToText: true,
        animateLLM: false,
        animateTool: false,
        animateTextToSpeech: false,
      });
      await new Promise(res => setTimeout(res, 4000));

      // Step 3: LLM
      setSteps({
        animateChannel: false,
        animateSpeechToText: false,
        animateLLM: true,
        animateTool: false,
        animateTextToSpeech: false,
      });
      await new Promise(res => setTimeout(res, 3000));

      // Step 4: Tool
      setSteps({
        animateChannel: false,
        animateSpeechToText: false,
        animateLLM: false,
        animateTool: true,
        animateTextToSpeech: false,
      });
      await new Promise(res => setTimeout(res, 4000));

      // Step 5: TextToSpeech
      setSteps({
        animateChannel: false,
        animateSpeechToText: false,
        animateLLM: false,
        animateTool: false,
        animateTextToSpeech: true,
      });
      await new Promise(res => setTimeout(res, 3000));

      // Reset idle state
      setSteps({
        animateChannel: false,
        animateSpeechToText: false,
        animateLLM: false,
        animateTool: false,
        animateTextToSpeech: false,
      });

      await new Promise(res => setTimeout(res, 1000)); // Optional: Add buffer before next loop
    };

    const startSequence = async () => {
      while (true) {
        await sequence();
      }
    };

    startSequence();

    return () => {
      // No interval to clear since we're using a loop
    };
  }, []);
  return (
    <div className="border sm:rounded-2xl bg-white dark:bg-gray-950">
      <Channel isAnimate={steps.animateChannel} />
      <SpeechToText isAnimate={steps.animateSpeechToText} />
      <LLM isAnimate={steps.animateLLM} />
      <Tool isAnimate={steps.animateTool} />
      <TextToSpeech isAnimate={steps.animateTextToSpeech} />
    </div>
  );
};

export const Channel: FC<{ isAnimate: boolean }> = (props: { isAnimate }) => {
  return (
    <motion.div
      className={cn(
        'sm:rounded-t-2xl',
        'transition-all delay-300',
        'w-fit z-10 animate-border-rotate [background:linear-gradient(45deg,#ffffff,--theme(--color-white)_50%,#ffffff)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-200/.48)_80%,--theme(--color-blue-500)_86%,--theme(--color-blue-300)_90%,--theme(--color-blue-500)_94%,--theme(--color-gray-600/.48))_border-box] border-[2px] border-transparent',
        'dark:[background:linear-gradient(45deg,#000000,--theme(--color-black)_50%,#000000)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-800/.48)_80%,--theme(--color-blue-500)_86%,--theme(--color-blue-300)_90%,--theme(--color-blue-500)_94%,--theme(--color-gray-800/.48))_border-box]',
        '@container relative flex h-full items-center justify-center w-full',
        !props.isAnimate
          ? '[background:transparent]! border-transparent! opacity-50'
          : 'opacity-100',
      )}
    >
      <div className="flex gap-3 relative">
        <div className="flex items-center gap-3 py-10 text-sm/6 font-semibold text-gray-950 dark:text-white border-x dark:border-gray-800/50 px-6 ">
          <Phone
            className="size-6 flex-none text-gray-600 dark:text-gray-400"
            strokeWidth={1.5}
          />
          Phone
        </div>
        <div className="flex items-center gap-3 text-sm/6 font-semibold text-gray-950 dark:text-white border-r dark:border-gray-800/50 px-6 ">
          <Globe2
            className="size-6 flex-none text-gray-600 dark:text-gray-400"
            strokeWidth={1.5}
          />
          Web
        </div>
        <div className="flex items-center gap-3 text-sm/6 font-semibold text-gray-950 dark:text-white border-r dark:border-gray-800/50 px-6 ">
          <TabletSmartphone
            className="size-6 flex-none text-gray-600 dark:text-gray-400"
            strokeWidth={1.5}
          />
          App
        </div>
        <div className="flex items-center gap-3 text-sm/6 font-semibold text-gray-950 dark:text-white border-r dark:border-gray-800/50 px-6 ">
          <svg
            className="size-6 flex-none fill-gray-600 dark:fill-gray-400"
            viewBox="0 0 32 32"
            version="1.1"
          >
            <path
              stroke="none"
              d="M26.576 5.363c-2.69-2.69-6.406-4.354-10.511-4.354-8.209 0-14.865 6.655-14.865 14.865 0 2.732 0.737 5.291 2.022 7.491l-0.038-0.070-2.109 7.702 7.879-2.067c2.051 1.139 4.498 1.809 7.102 1.809h0.006c8.209-0.003 14.862-6.659 14.862-14.868 0-4.103-1.662-7.817-4.349-10.507l0 0zM16.062 28.228h-0.005c-0 0-0.001 0-0.001 0-2.319 0-4.489-0.64-6.342-1.753l0.056 0.031-0.451-0.267-4.675 1.227 1.247-4.559-0.294-0.467c-1.185-1.862-1.889-4.131-1.889-6.565 0-6.822 5.531-12.353 12.353-12.353s12.353 5.531 12.353 12.353c0 6.822-5.53 12.353-12.353 12.353h-0zM22.838 18.977c-0.371-0.186-2.197-1.083-2.537-1.208-0.341-0.124-0.589-0.185-0.837 0.187-0.246 0.371-0.958 1.207-1.175 1.455-0.216 0.249-0.434 0.279-0.805 0.094-1.15-0.466-2.138-1.087-2.997-1.852l0.010 0.009c-0.799-0.74-1.484-1.587-2.037-2.521l-0.028-0.052c-0.216-0.371-0.023-0.572 0.162-0.757 0.167-0.166 0.372-0.434 0.557-0.65 0.146-0.179 0.271-0.384 0.366-0.604l0.006-0.017c0.043-0.087 0.068-0.188 0.068-0.296 0-0.131-0.037-0.253-0.101-0.357l0.002 0.003c-0.094-0.186-0.836-2.014-1.145-2.758-0.302-0.724-0.609-0.625-0.836-0.637-0.216-0.010-0.464-0.012-0.712-0.012-0.395 0.010-0.746 0.188-0.988 0.463l-0.001 0.002c-0.802 0.761-1.3 1.834-1.3 3.023 0 0.026 0 0.053 0.001 0.079l-0-0.004c0.131 1.467 0.681 2.784 1.527 3.857l-0.012-0.015c1.604 2.379 3.742 4.282 6.251 5.564l0.094 0.043c0.548 0.248 1.25 0.513 1.968 0.74l0.149 0.041c0.442 0.14 0.951 0.221 1.479 0.221 0.303 0 0.601-0.027 0.889-0.078l-0.031 0.004c1.069-0.223 1.956-0.868 2.497-1.749l0.009-0.017c0.165-0.366 0.261-0.793 0.261-1.242 0-0.185-0.016-0.366-0.047-0.542l0.003 0.019c-0.092-0.155-0.34-0.247-0.712-0.434z"
            ></path>
          </svg>
          WhatsApp
        </div>
      </div>
    </motion.div>
  );
};

export const SpeechToText: FC<{ isAnimate: boolean }> = (props: {
  isAnimate;
}) => {
  return (
    <motion.div
      className={cn(
        'transition-all delay-300',
        'w-fit z-10 animate-border-rotate [background:linear-gradient(45deg,#ffffff,--theme(--color-white)_50%,#ffffff)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-200/.48)_80%,--theme(--color-amber-500)_86%,--theme(--color-amber-300)_90%,--theme(--color-amber-500)_94%,--theme(--color-gray-600/.48))_border-box] border-[2px] border-transparent',
        'dark:[background:linear-gradient(45deg,#000000,--theme(--color-black)_50%,#000000)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-800/.48)_80%,--theme(--color-amber-500)_86%,--theme(--color-amber-300)_90%,--theme(--color-amber-500)_94%,--theme(--color-gray-800/.48))_border-box]',
        '@container relative flex flex-col h-full items-center justify-center w-full',
        'block overflow-hidden relative',
        !props.isAnimate
          ? '[background:transparent]! border-transparent! opacity-50'
          : 'opacity-100',
      )}
    >
      <div className=" z-10 border-y border-amber-600/50 bg-amber-500/10 text-amber-600">
        <div className="py-2 text-center font-mono text-sm font-semibold tracking-widest uppercase opacity-60">
          Speech to Text
        </div>
      </div>
      {/*  */}
      <div className="@container relative flex items-center justify-center w-full ">
        <motion.div
          //   transition={{
          //     duration: 10,
          //     ease: 'linear',
          //     repeat: Infinity,
          //   }}
          //   initial={{ translateX: 0 }}
          //   animate={props.isAnimate ? { translateX: '-10%' } : {}}
          className="flex relative h-full divide-x dark:divide-gray-800/50 no-scrollbar "
        >
          <div className="flex items-center gap-3 py-10 text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <svg
              className="size-6"
              viewBox="0 0 96 96"
              xmlns="http://www.w3.org/2000/svg"
            >
              <defs>
                <linearGradient
                  id="e399c19f-b68f-429d-b176-18c2117ff73c"
                  x1="-1032.172"
                  x2="-1059.213"
                  y1="145.312"
                  y2="65.426"
                  gradientTransform="matrix(1 0 0 -1 1075 158)"
                  gradientUnits="userSpaceOnUse"
                >
                  <stop offset="0" stopColor="#114a8b" />
                  <stop offset="1" stopColor="#0669bc" />
                </linearGradient>
                <linearGradient
                  id="ac2a6fc2-ca48-4327-9a3c-d4dcc3256e15"
                  x1="-1023.725"
                  x2="-1029.98"
                  y1="108.083"
                  y2="105.968"
                  gradientTransform="matrix(1 0 0 -1 1075 158)"
                  gradientUnits="userSpaceOnUse"
                >
                  <stop offset="0" stopOpacity=".3" />
                  <stop offset=".071" stopOpacity=".2" />
                  <stop offset=".321" stopOpacity=".1" />
                  <stop offset=".623" stopOpacity=".05" />
                  <stop offset="1" stopOpacity="0" />
                </linearGradient>
                <linearGradient
                  id="a7fee970-a784-4bb1-af8d-63d18e5f7db9"
                  x1="-1027.165"
                  x2="-997.482"
                  y1="147.642"
                  y2="68.561"
                  gradientTransform="matrix(1 0 0 -1 1075 158)"
                  gradientUnits="userSpaceOnUse"
                >
                  <stop offset="0" stopColor="#3ccbf4" />
                  <stop offset="1" stopColor="#2892df" />
                </linearGradient>
              </defs>
              <path
                fill="url(#e399c19f-b68f-429d-b176-18c2117ff73c)"
                d="M33.338 6.544h26.038l-27.03 80.087a4.152 4.152 0 0 1-3.933 2.824H8.149a4.145 4.145 0 0 1-3.928-5.47L29.404 9.368a4.152 4.152 0 0 1 3.934-2.825z"
              />
              <path
                fill="#0078d4"
                d="M71.175 60.261h-41.29a1.911 1.911 0 0 0-1.305 3.309l26.532 24.764a4.171 4.171 0 0 0 2.846 1.121h23.38z"
              />
              <path
                fill="url(#ac2a6fc2-ca48-4327-9a3c-d4dcc3256e15)"
                d="M33.338 6.544a4.118 4.118 0 0 0-3.943 2.879L4.252 83.917a4.14 4.14 0 0 0 3.908 5.538h20.787a4.443 4.443 0 0 0 3.41-2.9l5.014-14.777 17.91 16.705a4.237 4.237 0 0 0 2.666.972H81.24L71.024 60.261l-29.781.007L59.47 6.544z"
              />
              <path
                fill="url(#a7fee970-a784-4bb1-af8d-63d18e5f7db9)"
                d="M66.595 9.364a4.145 4.145 0 0 0-3.928-2.82H33.648a4.146 4.146 0 0 1 3.928 2.82l25.184 74.62a4.146 4.146 0 0 1-3.928 5.472h29.02a4.146 4.146 0 0 0 3.927-5.472z"
              />
            </svg>
            Azure
          </div>
          <div className="flex items-center gap-3 py-10 text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/198796716894742118.jpg"
            />
            Google
          </div>
          <div className="flex items-center gap-3 py-10 text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/2123891723608588082.jpg"
            />
            Deepgram
          </div>
          <div className="flex items-center gap-3 py-10 text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/assemblyai.png"
            />
            AssemblyAI
          </div>
          <div className="flex items-center gap-3 py-10 text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://avatars.githubusercontent.com/u/92447723?s=200&v=4"
            />
            Gladia
          </div>
          <div className="flex items-center gap-3 py-10 text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://upload.wikimedia.org/wikipedia/commons/1/1a/SM-Icon-Dark_Cyan1000.png"
            />
            Speechmatics
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSobT6Nq7W-FJnK5lLapZlwySLwB0W4sKCYDg&s"
            />
            OpenAI Whisper
          </div>
        </motion.div>
      </div>
    </motion.div>
  );
};

export const TextToSpeech: FC<{ isAnimate: boolean }> = (props: {
  isAnimate;
}) => {
  return (
    <motion.div
      className={cn(
        'sm:rounded-bl-2xl transition-all delay-300',
        'w-fit z-10 animate-border-rotate [background:linear-gradient(45deg,#ffffff,--theme(--color-white)_50%,#ffffff)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-200/.48)_80%,--theme(--color-purple-500)_86%,--theme(--color-rose-300)_90%,--theme(--color-rose-500)_94%,--theme(--color-gray-600/.48))_border-box] border-[2px] border-transparent',
        'dark:[background:linear-gradient(45deg,#000000,--theme(--color-black)_50%,#000000)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-800/.48)_80%,--theme(--color-rose-500)_86%,--theme(--color-rose-300)_90%,--theme(--color-rose-500)_94%,--theme(--color-gray-800/.48))_border-box]',
        '@container relative flex flex-col h-full items-center justify-center w-full',
        'block overflow-hidden relative',
        !props.isAnimate
          ? '[background:transparent]! border-transparent! opacity-50'
          : 'opacity-100',
      )}
    >
      <div className=" z-10 border-y border-rose-600/50 bg-rose-500/10 text-rose-600">
        <div className="py-2 text-center font-mono text-sm font-semibold tracking-widest uppercase opacity-60">
          Text To Speech
        </div>
      </div>
      <div className="@container relative flex items-center justify-center w-full ">
        <motion.div
          transition={{
            duration: 10,
            ease: 'linear',
            repeat: Infinity,
          }}
          initial={{ translateX: 0 }}
          animate={props.isAnimate ? { translateX: '-10%' } : {}}
          className="flex relative h-full divide-x no-scrollbar"
        >
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/cartesia.jpg"
              alt="Cartesia"
            />
            Cartesia
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/11labs.png"
              alt="ElevenLabs"
            />
            ElevenLabs
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/198796716894742122.png"
              alt="Azure"
            />
            Azure
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/198796716894742118.jpg"
              alt="Google"
            />
            Google
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/2123891723608588082.jpg"
              alt="Deepgram"
            />
            Deepgram
          </div>
        </motion.div>
      </div>
    </motion.div>
  );
};

export const LLM: FC<{ isAnimate: boolean }> = (props: { isAnimate }) => {
  return (
    <motion.div
      className={cn(
        'transition-all delay-300',
        'w-fit z-10 animate-border-rotate [background:linear-gradient(45deg,#ffffff,--theme(--color-white)_50%,#ffffff)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-200/.48)_80%,--theme(--color-purple-500)_86%,--theme(--color-purple-300)_90%,--theme(--color-purple-500)_94%,--theme(--color-gray-600/.48))_border-box] border-[2px] border-transparent',
        'dark:[background:linear-gradient(45deg,#000000,--theme(--color-black)_50%,#000000)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-800/.48)_80%,--theme(--color-purple-500)_86%,--theme(--color-purple-300)_90%,--theme(--color-purple-500)_94%,--theme(--color-gray-800/.48))_border-box]',
        '@container relative flex flex-col h-full items-center justify-center w-full',
        'block overflow-hidden relative',
        !props.isAnimate
          ? '[background:transparent]! border-transparent! opacity-50'
          : 'opacity-100',
      )}
    >
      <div className=" z-10 border-y border-purple-600/50 bg-purple-500/10 text-purple-600">
        <div className="py-2 text-center font-mono text-sm font-semibold tracking-widest uppercase opacity-60">
          LLM
        </div>
      </div>
      <div className="@container relative flex h-full items-center justify-center w-full ">
        <motion.div
          transition={{
            duration: 10,
            ease: 'linear',
            repeat: Infinity,
          }}
          initial={{ translateX: 0 }}
          animate={!props.isAnimate ? { translateX: '-10%' } : {}}
          className="flex relative h-full divide-x no-scrollbar dark:divide-gray-800/50"
        >
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/198796716894742122.png"
            />
            Azure
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/1987967168347635712.jpg"
            />
            Anthropic
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/1987967168435716096.png"
            />
            Cohere
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRjY2Oj_qhiyxddGNn9BLCKA-cf_M5d87kUXA&s"
            />
            Deepseek
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-8 flex-none text-gray-600 dark:text-gray-400"
              src="https://miro.medium.com/1*-U0stHO5R9cxKUqkJL-9-w.jpeg"
            />
            Cerebras
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-8 flex-none text-gray-600 dark:text-gray-400"
              src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQnuyhgsRh7ORw2Lh_SAh0XXkI6R9YvvaDCxg&s"
            />
            Gork
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/1987967168452493312.svg"
            />
            OpenAI
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://rapida-assets-01.s3.ap-south-1.amazonaws.com/providers/198796716894742118.jpg"
            />
            Google
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Workflow
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Workflows
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Code
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Websocket
          </div>
          <div className="flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Agentkit
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6 w-max">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Marketing and Sales Agent
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Collection Agent
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Customer Support
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Onboarding Agent
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Front Desk Agent
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Reminders Agent
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Lead Qualification Agent
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Bot
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Surveys Agent
          </div>
        </motion.div>
      </div>
    </motion.div>
  );
};

const Tool: FC<{ isAnimate: boolean }> = (props: { isAnimate }) => {
  return (
    <motion.div
      className={cn(
        'col-span-3',
        'transition-all delay-300',
        'w-fit z-10 animate-border-rotate [background:linear-gradient(45deg,#ffffff,--theme(--color-white)_50%,#ffffff)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-200/.48)_80%,--theme(--color-teal-500)_86%,--theme(--color-teal-300)_90%,--theme(--color-teal-500)_94%,--theme(--color-gray-600/.48))_border-box] border-[2px] border-transparent',
        'dark:[background:linear-gradient(45deg,#000000,--theme(--color-black)_50%,#000000)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-gray-800/.48)_80%,--theme(--color-teal-500)_86%,--theme(--color-teal-300)_90%,--theme(--color-teal-500)_94%,--theme(--color-gray-800/.48))_border-box]',
        '@container relative flex flex-col h-full items-center justify-center w-full',
        'block overflow-hidden relative',
        !props.isAnimate
          ? '[background:transparent]! border-transparent! opacity-50'
          : 'opacity-100',
      )}
    >
      <div className=" z-10 border-y border-teal-600/50 bg-teal-500/10 text-teal-600">
        <div className="py-2 text-center font-mono text-sm font-semibold tracking-widest uppercase opacity-60">
          Prebuild tools
        </div>
      </div>
      <div className="  @container relative flex h-full items-center justify-center w-full">
        <motion.div
          transition={{
            duration: 10,
            ease: 'linear',
            repeat: Infinity,
          }}
          initial={{ translateX: 0 }}
          animate={props.isAnimate ? { translateX: '-10%' } : {}}
          className="flex relative h-full divide-x no-scrollbar dark:divide-gray-800/50"
        >
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <PhoneForwarded
              className="w-6 h-6 text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Transfer Call
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <MessageCircle
              className="w-6 h-6 text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Route Agent
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <BookOpen
              className="w-6 h-6 text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Knowledge Retrieval
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Cloud
              className="w-6 h-6 text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            API Request
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <RefreshCcw
              className="w-6 h-6 text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Warm Transfer
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <Mail
              className="w-6 h-6 text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Send Email
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6">
            <MessageSquare
              className="w-6 h-6 text-gray-600 dark:text-gray-400"
              strokeWidth={1.5}
            />
            Send SMS
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6 w-max">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://i.pinimg.com/736x/fa/83/b8/fa83b8de8ff5a8e2dc98c6f7e249c842.jpg"
            />
            MCP
          </div>
          <div className="whitespace-nowrap flex items-center gap-3 py-10  text-sm/6 font-semibold text-gray-950 dark:text-white px-6 w-max">
            <img
              className="size-6 flex-none text-gray-600 dark:text-gray-400"
              src="https://www.pipelinersales.com/wp-content/uploads/2018/07/zapier.jpg"
            />
            Zapier
          </div>
        </motion.div>
      </div>
    </motion.div>
  );
};
