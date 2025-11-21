import { cn } from '@/utils';

/**
 *  */
interface InputHelperProp extends React.LabelHTMLAttributes<HTMLLabelElement> {}

/**
 *
 * @param props
 * @returns
 */
export function InputHelper(props: InputHelperProp) {
  return (
    <div
      className={cn(
        'mb-3 text-sm/6 dark:text-gray-500 text-gray-500',
        props.className,
      )}
    >
      {props.children}
    </div>
  );
}
