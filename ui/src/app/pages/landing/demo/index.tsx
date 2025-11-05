import { FormLabel } from '@/app/components/form-label';
import { ErrorMessage } from '@/app/components/Form/error-message';
import { FieldSet } from '@/app/components/Form/Fieldset';
import { Input } from '@/app/components/Form/Input';
import { Select } from '@/app/components/Form/Select';
import { SuccessMessage } from '@/app/components/Form/success-message';
import { Helmet } from '@/app/components/Helmet';
import { InputHelper } from '@/app/components/input-helper';
import { Footer } from '@/app/pages/landing/components/footer';
import { Header } from '@/app/pages/landing/components/header';
import { connectionConfig } from '@/configs';
import { CreateLead, CreateLeadRequest } from '@rapidaai/react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

export const LeadGeneration = () => {
  //
  const { register, handleSubmit } = useForm();
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  /**
   * calling for creating a lead
   * @param data
   */
  const onCreateLead = async data => {
    console.dir(data);
    setError('');
    if (!data.company) {
      setError('Please provide a valid company name and try again.');
      return;
    }
    if (!data.email) {
      setError('Please provide a valid business email and try again.');
      return;
    }
    const request = new CreateLeadRequest();
    request.setCompany(data.company);
    request.setEmail(data.email);
    request.setExpectedvolume(data.expectedVolume);
    CreateLead(connectionConfig, request)
      .then(x => {
        if (x.getSuccess()) {
          setSuccess(
            'Thank you, Our team will reach out shortly to schedule a session with expert.',
          );
          return;
        }
        if (x.getError()?.getHumanmessage()) {
          setError(x.getError()?.getHumanmessage()!);
        }
      })
      .catch(err => {
        setError('Unable to complete request, please try again in sometime.');
      });
  };

  return (
    <>
      <Helmet title="Meet with an voice AI expert"></Helmet>
      <div
        className="font-sans text-gray-950 dark:text-white antialiased [overflow-anchor:none] container mx-auto"
        scroll-region=""
      >
        <div className="isolate">
          <Header />
          <main className="flex min-h-dvh flex-col pt-14">
            <div className="grid flex-1 grid-rows-[1fr_auto] overflow-clip grid-cols-[1fr_var(--gutter-width)_minmax(0,var(--breakpoint-2xl))_var(--gutter-width)_1fr] [--gutter-width:--spacing(6)] lg:[--gutter-width:--spacing(10)]">
              <div className="col-start-2 row-span-full row-start-1 max-sm:hidden border-x  bg-size-[10px_10px] bg-fixed bg-[repeating-linear-gradient(315deg,theme('colors.gray.200')_0,theme('colors.gray.200')_1px,transparent_0,transparent_50%)] dark:bg-[repeating-linear-gradient(315deg,theme('colors.gray.800')_0,theme('colors.gray.800')_1px,transparent_0,transparent_50%)]" />
              <div className="col-start-4 row-span-full row-start-1 max-sm:hidden border-x  bg-size-[10px_10px] bg-fixed bg-[repeating-linear-gradient(315deg,theme('colors.gray.200')_0,theme('colors.gray.200')_1px,transparent_0,transparent_50%)] dark:bg-[repeating-linear-gradient(315deg,theme('colors.gray.800')_0,theme('colors.gray.800')_1px,transparent_0,transparent_50%)]" />

              {/* hero */}
              <div className="col-start-3 row-start-1 max-sm:col-span-full max-sm:col-start-1">
                <div className="border-y mt-12 grid gap-x-10 sm:mt-20 lg:mt-24 lg:grid-cols-[1fr_1fr]">
                  <div className="py-2 max-lg:line-b lg:border-r">
                    <h1 className="px-4 sm:px-2 mt-2 text-3xl sm:text-4xl text-pretty">
                      See how leading enterprises scale Voice AI with
                      orchestration
                    </h1>
                    <p className="px-4 mt-4 sm:px-2">
                      Book a free consultation with our Voice AI experts to
                      explore how you can deliver natural, real-time
                      conversations that improve customer experience and reduce
                      operational costs.
                    </p>

                    <div className="mt-10 relative border-y">
                      <p className="px-4 text-lg sm:px-2">What’s included</p>
                    </div>

                    <ul className="px-4 sm:px-2 py-4 group grid grid-cols-1 gap-x-6 gap-y-2 text-sm/7 text-gray-600 data-dark:text-gray-300 @3xl:grid-cols-2 dark:text-gray-400">
                      <li className="flex gap-2">
                        <svg
                          aria-hidden="true"
                          viewBox="0 0 22 22"
                          className="h-7 w-5.5"
                        >
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M22 11c0 6.075-4.925 11-11 11S0 17.075 0 11 4.925 0 11 0s11 4.925 11 11Z"
                          ></path>
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M11 21c5.523 0 10-4.477 10-10S16.523 1 11 1 1 5.477 1 11s4.477 10 10 10Zm0 1c6.075 0 11-4.925 11-11S17.075 0 11 0 0 4.925 0 11s4.925 11 11 11Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                          <path
                            className="fill-white"
                            d="m14.684 7.82-4.079 6.992L7.293 11.5 8 10.793l2.395 2.395 3.425-5.872.864.504Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                        </svg>
                        <p>
                          A review of your current voice and call automation
                          setup.
                        </p>
                      </li>
                      <li className="flex gap-2">
                        <svg
                          aria-hidden="true"
                          viewBox="0 0 22 22"
                          className="h-7 w-5.5"
                        >
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M22 11c0 6.075-4.925 11-11 11S0 17.075 0 11 4.925 0 11 0s11 4.925 11 11Z"
                          ></path>
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M11 21c5.523 0 10-4.477 10-10S16.523 1 11 1 1 5.477 1 11s4.477 10 10 10Zm0 1c6.075 0 11-4.925 11-11S17.075 0 11 0 0 4.925 0 11s4.925 11 11 11Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                          <path
                            className="fill-white"
                            d="m14.684 7.82-4.079 6.992L7.293 11.5 8 10.793l2.395 2.395 3.425-5.872.864.504Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                        </svg>
                        <p>
                          Tailored suggestions for using Voice AI across your
                          channels.
                        </p>
                      </li>
                      <li className="flex gap-2">
                        <svg
                          aria-hidden="true"
                          viewBox="0 0 22 22"
                          className="h-7 w-5.5"
                        >
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M22 11c0 6.075-4.925 11-11 11S0 17.075 0 11 4.925 0 11 0s11 4.925 11 11Z"
                          ></path>
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M11 21c5.523 0 10-4.477 10-10S16.523 1 11 1 1 5.477 1 11s4.477 10 10 10Zm0 1c6.075 0 11-4.925 11-11S17.075 0 11 0 0 4.925 0 11s4.925 11 11 11Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                          <path
                            className="fill-white"
                            d="m14.684 7.82-4.079 6.992L7.293 11.5 8 10.793l2.395 2.395 3.425-5.872.864.504Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                        </svg>
                        <p>A custom demo of Rapida’s orchestration platform</p>
                      </li>
                      <li className="flex gap-2">
                        <svg
                          aria-hidden="true"
                          viewBox="0 0 22 22"
                          className="h-7 w-5.5"
                        >
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M22 11c0 6.075-4.925 11-11 11S0 17.075 0 11 4.925 0 11 0s11 4.925 11 11Z"
                          ></path>
                          <path
                            className="fill-gray-500 dark:fill-gray-700"
                            d="M11 21c5.523 0 10-4.477 10-10S16.523 1 11 1 1 5.477 1 11s4.477 10 10 10Zm0 1c6.075 0 11-4.925 11-11S17.075 0 11 0 0 4.925 0 11s4.925 11 11 11Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                          <path
                            className="fill-white"
                            d="m14.684 7.82-4.079 6.992L7.293 11.5 8 10.793l2.395 2.395 3.425-5.872.864.504Z"
                            clipRule="evenodd"
                            fillRule="evenodd"
                          ></path>
                        </svg>
                        <p>
                          Real examples of enterprises scaling securely with
                          Rapida.
                        </p>
                      </li>
                    </ul>
                  </div>
                  <form
                    method="POST"
                    className="grid grid-cols-1 grid-rows-[1fr_auto] lg:border-l"
                    onSubmit={handleSubmit(onCreateLead)}
                  >
                    <div className="px-4 sm:px-2 py-6 space-y-6">
                      <div className="text-xl">
                        Book a quick session with our team to understand how
                        Rapida can fit into your existing stack and help you
                        scale voice AI confidently.
                      </div>
                      <div className="flex w-full space-x-4">
                        <FieldSet className="w-full">
                          <FormLabel>Business Email*</FormLabel>
                          <Input
                            autoComplete="name"
                            type="email"
                            required
                            className="bg-light-background"
                            placeholder="eg: john@rapida.ai"
                            {...register('email')}
                          ></Input>
                        </FieldSet>
                        <FieldSet className="w-full">
                          <FormLabel>Company*</FormLabel>
                          <Input
                            autoComplete="name"
                            type="text"
                            required
                            className="bg-light-background"
                            placeholder="eg: John deo"
                            {...register('company')}
                          ></Input>
                        </FieldSet>
                      </div>
                      <FieldSet className="w-full">
                        <FormLabel className="normal-case">
                          What is your expected volume this year across all
                          communication channels?*
                        </FormLabel>
                        <Select
                          className="bg-light-background"
                          {...register('expectedVolume')}
                          options={[
                            {
                              name: "I don't know",
                              value: "I don't know",
                            },
                            {
                              name: '0 - 99,999',
                              value: '0 - 99,999',
                            },
                            {
                              name: '100,000 - 299,999',
                              value: '100,000 - 299,999',
                            },
                            {
                              name: '300,000 - 499,999',
                              value: '300,000 - 499,999',
                            },
                            {
                              name: '500,000 - 999,999',
                              value: '500,000 - 999,999',
                            },
                            {
                              name: '1M - 4.99M',
                              value: '1M - 4.99M',
                            },
                            {
                              name: '5M - 19.99M',
                              value: '5M - 19.99M',
                            },
                            {
                              name: '20M - 99.99M',
                              value: '20M - 99.99M',
                            },
                            {
                              name: 'More than 100 million',
                              value: 'More than 100 million',
                            },
                          ]}
                        ></Select>
                      </FieldSet>
                      <InputHelper>
                        By clicking submit, you acknowledge our
                        <a
                          className="underline hover:text-blue-600 mx-2"
                          href="/static/privacy-policy"
                        >
                          Privacy Policy
                        </a>
                        and agree to receive email communications from us. You
                        can unsubscribe at any time.
                      </InputHelper>
                      <ErrorMessage message={error} />
                      <SuccessMessage message={success} />
                    </div>
                    <div className="flex gap-4 px-4 sm:px-2 py-2 whitespace-nowrap border-t">
                      <button
                        className="pl-4 px-2 py-2 gap-2 inline-flex justify-center rounded-full text-base font-medium focus-visible:outline-2 focus-visible:outline-offset-2 bg-blue-600 text-white cursor-pointer hover:bg-blue-600 focus:outline-2 focus:outline-offset-2 focus:outline-blue-500 active:bg-blue-700"
                        type="submit"
                      >
                        Consult an expert
                        <span className="bg-blue-800 rounded-full p-2 w-6 h-6 items-center flex">
                          <svg
                            fill="currentColor"
                            aria-hidden="true"
                            viewBox="0 0 10 10"
                            className="-mr-0.5 w-4"
                          >
                            <path d="M4.85355 0.146423L9.70711 4.99998L4.85355 9.85353L4.14645 9.14642L7.79289 5.49998H0V4.49998H7.79289L4.14645 0.85353L4.85355 0.146423Z" />
                          </svg>
                        </span>
                      </button>
                    </div>
                  </form>
                </div>
                <section className="mt-20 sm:mt-40">
                  <div className="border-y px-4 py-2 sm:px-2">
                    <h2 className="max-w-3xl text-3xl sm:text-4xl text-pretty font-medium">
                      Results from Voice AI first automation
                    </h2>
                  </div>
                  <div className="bg-gray-950/5 dark:bg-white/5 p-4 grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    {/*  */}
                    <div className="grid grid-cols-1 rounded-2xl bg-purple-500/50 text-white aspect-square relative">
                      <span className="text-4xl p-6">40%</span>
                      <p className="text-lg absolute bottom-0 p-6">
                        Reduce average handling time across inbound calls with
                        intelligent voice orchestration
                      </p>
                    </div>
                    <div className="grid grid-cols-1 rounded-2xl relative bg-green-500/50 text-white aspect-square">
                      <span className="text-4xl p-6">60%</span>
                      <p className="text-lg absolute bottom-0 p-6">
                        Lower cost per resolution while keeping full control and
                        compliance
                      </p>
                    </div>
                    <div className="grid grid-cols-1 rounded-2xl relative bg-sky-500/30 text-white aspect-square">
                      <span className="text-4xl p-6">85%</span>
                      <p className="text-lg absolute bottom-0 p-6">
                        Handle end-to-end conversations naturally through ASR +
                        Agent + TTS
                      </p>
                    </div>
                    <div className="grid grid-cols-1 rounded-2xl relative bg-sky-500/30 text-white aspect-square">
                      <span className="text-4xl p-6">3x</span>
                      <p className="text-lg absolute bottom-0 p-6">
                        Deliver human-like, context-aware responses that
                        eliminate IVR frustration.
                      </p>
                    </div>
                  </div>
                  {/*  */}
                </section>
              </div>
              <Footer />
            </div>
          </main>
        </div>
      </div>
    </>
  );
};
