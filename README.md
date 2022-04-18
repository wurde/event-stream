# Event Stream

An example event-driven architecture using NATS JetStream.

This app consumes several RSS and Atom feeds to produce
a daily email digest with the latest entries.

Feeds included (but customizable given basic coding skills):

    SEC 10-K filings

    https://aws.amazon.com/blogs/machine-learning/feed
    
    https://news.mit.edu/topic/mitmachine-learning-rss.xml
    
    http://googleaiblog.blogspot.com/atom.xml
    
    https://openai.com/blog/rss
    
    https://research.facebook.com/feed
    
    https://developer.nvidia.com/blog/feed
    
    https://bair.berkeley.edu/blog/feed.xml

## Getting Started

### 1. Set Environment Variables

Set the AWS access keys. Access keys consist of an access key ID
and secret access key, which are used to sign programmatic
requests that you make to AWS. If you don’t have access keys,
you can create them by using the AWS Management Console. We
recommend that you use IAM access keys instead of AWS root
account access keys. IAM lets you securely control access to AWS
services and resources in your AWS account.

Also set the email you want to receive the daily digests.

It supports setting environment variables via Dotenv (`.env`).

```env
<!-- REQUIRED -->
AWS_ACCESS_KEY_ID=<CV2XP27...>
AWS_SECRET_ACCESS_KEY=<gLcGMJ1z...>
FROM_EMAIL=<hello@example.com>
TO_EMAIL=<hello@example.com>
```

### 3. Verify Email Address

Before you can send an email using Amazon SES, you must
create and verify each identity that you're going to use.
You likely have to check the spam folder the first time to
mark it as not spam.

### 2. Build and Run the Application:

```bash
git clone https://github.com/wurde/event-stream
cd event-stream

go build

./event-stream
```

## How It Works

Emails are sent via AWS SES. Amazon SES is an email platform
that provides an easy, cost-effective way for you to send
and receive email. Pay as you go pricing:

    $0 for the first 62,000 emails you send each month, and
    $0.10 for every 1,000 emails you send after that.
    See latest https://aws.amazon.com/ses/pricing

## License

This project is __FREE__ to use, reuse, remix, and resell.
This is made possible by the [MIT license](/LICENSE).
